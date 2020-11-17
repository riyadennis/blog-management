package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/riyadennis/blog-management/events"
	"os"
)

// Article holds structure of the blob in Article event
type Article struct {
	CreatedBy string
	Name      string
	Heading   string
	Body      string
}

type EventStore interface {
	Add(ctx context.Context, e *events.Article) error
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

func (c *Config) Add(ctx context.Context, e *events.Article) error {
	query, err := c.Conn.Prepare("INSERT INTO events_store(id,version,state,data) values(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.Exec(e.ID, e.EventVersion, e.State, e.Content)
	if err != nil {
		return err
	}

	return nil
}
