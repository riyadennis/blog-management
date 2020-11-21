package eventsource

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/riyadennis/blog-management/pkg/api/events"
	"os"
	"time"
)

const TimeOut = 5

// Article holds structure of the blob in Article event
type Article struct {
	Author       string
	Heading      string
	Introduction string
	Body         string
}

type EventStore interface {
	Apply(ctx context.Context, e events.Event) error
	LatestVersion(ctx context.Context, aggregateId string) (int64, error)
	Load(ctx context.Context, aggregateId string) ([]events.Event, error)
}

type Config struct {
	Conn *sql.DB
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

	query, err := c.Conn.Prepare("INSERT INTO events_store(resourceID,version,state,content) values(?,?,?,?)")
	if err != nil {
		return err
	}

	ac, err := json.Marshal(model.Content)
	if err != nil {
		return err
	}

	_, err = query.ExecContext(ctx, model.ID, model.Version, model.State, string(ac))
	if err != nil {
		return err
	}

	return nil
}
