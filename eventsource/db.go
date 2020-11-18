package eventsource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/riyadennis/blog-management/events"
	"os"
	"time"
)

const TimeOut = 5

// Article holds structure of the blob in Article event
type Article struct {
	Author       string
	Heading      string
	Introduction string
	End          string
}

type EventStore interface {
	Apply(ctx context.Context, e events.Event) error
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

	query, err := c.Conn.Prepare("INSERT INTO events_store(id,version,state,data) values(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.ExecContext(ctx, model.ID, model.Version, model.State, model.Content)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Load(ctx context.Context, aggregateID string) ([]events.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)
	defer cancel()

	if aggregateID == "" {
		return nil, errors.New("empty aggregate id")
	}

	rows, err := c.Conn.QueryContext(ctx,
		"SELECT version,state,data,created_at FROM events_store WHERE id=?", aggregateID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var models []events.Event

	for rows.Next() {
		var state, data, createdAt string
		var version int

		if err := rows.Scan(&version, &state, &data, &createdAt); err != nil {
			return nil, err
		}

		createdTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
		if err != nil {
			return nil, err
		}

		models = append(models, &events.Model{
			ID:        aggregateID,
			Version:   version,
			State:     state,
			Content:   data,
			CreatedAt: createdTime,
		})
	}

	return models, nil
}
