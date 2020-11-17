package db

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
	CreatedBy string
	Name      string
	Heading   string
	Body      string
}

type EventStore interface {
	Add(ctx context.Context, state string, e events.Event) error
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

func (c *Config) Add(ctx context.Context, state string, e events.Event) error {
	ctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)
	defer cancel()

	if e == nil {
		return errors.New("empty event")
	}
	query, err := c.Conn.Prepare("INSERT INTO events_store(id,version,state,data) values(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.ExecContext(ctx, e.AggregateID(), e.Version(), state, e.Data())
	if err != nil {
		return err
	}

	return nil
}
