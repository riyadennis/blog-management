package eventsource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/riyadennis/blog-management/pkg/api/events"
)

const TimeOut = 5

// Article holds structure of the blob in Article event
type Article struct {
	Author       string `json:"author,omitempty"`
	Heading      string `json:"heading,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	Body         string `json:"body,omitempty"`
}

var (
	config EventStore
	once   sync.Once
)

type Command interface {
	Apply(ctx context.Context, e events.Event) error
	Aggregate(ctx context.Context, aggregateID string) ([]byte, error)
}

type Query interface {
	Load(ctx context.Context, aggregateId string) ([]events.Event, error)
	Events(ctx context.Context) ([]*Article, error)
}

type EventStore interface {
	Command
	Query
}

type Config struct {
	Conn *sql.DB
}

func Set(c EventStore) {
	once.Do(func() {
		config = c
	})
}

func Get() EventStore {
	return config
}

func NewConn() (*Config, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Config{Conn: conn}, nil
}

func (c *Config) Apply(ctx context.Context, e events.Event) error {
	ctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)
	defer cancel()

	if e == nil {
		return errors.New("empty event")
	}

	model, ok := e.(*events.Model)
	if !ok {
		return errors.New("invalid event")
	}
	version, err := latestVersion(ctx, c.Conn, model.ID)
	if err != nil {
		return err
	}

	query, err := c.Conn.Prepare("INSERT INTO events_store(resourceID,version,state,content,aggregateID) values(?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.ExecContext(ctx, model.ID, version+1, model.State, model.Content, model.AggregateID)
	if err != nil {
		return err
	}

	return nil
}

func latestVersion(ctx context.Context, conn *sql.DB, resourceID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)
	defer cancel()

	if resourceID == "" {
		return 0, errors.New("empty aggregate id")
	}

	var version interface{}
	row := conn.QueryRowContext(ctx, "SELECT MAX(version) as version FROM events_store WHERE resourceID=?", resourceID)

	err := row.Scan(&version)
	if err != nil {
		return 0, err
	}

	if version == nil {
		return 0, nil
	}

	return version.(int64), nil
}
