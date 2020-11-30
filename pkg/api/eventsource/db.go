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

type EventStore interface {
	Apply(ctx context.Context, e events.Event) error
	LatestVersion(ctx context.Context, aggregateId string) (int64, error)
	Aggregate(ctx context.Context, aggregateID string) ([]byte, error)
	Load(ctx context.Context, aggregateId string) ([]events.Event, error)
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

	query, err := c.Conn.Prepare("INSERT INTO events_store(resourceID,version,state,content,aggregateID) values(?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.ExecContext(ctx, model.ID, model.Version, model.State, model.Content, model.AggregateID)
	if err != nil {
		return err
	}

	return nil
}
