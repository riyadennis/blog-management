package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// Article holds structure of the blob in Article
type Article struct {
	CreatedBy string
	Name      string
	Heading   string
	Body      string
}

type EventStore interface {
	Add(ctx context.Context, e Command) error
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

func (c *Config) Add(ctx context.Context, e Command) error {
	query, err := c.Conn.Prepare("INSERT INTO events_store(id,version,state,data) values(?,?,?,?)")
	if err != nil {
		return err
	}

	cc, ok := e.(CreateCommand)
	if !ok {
		return errors.New("invalid command")
	}

	_, err = query.Exec(cc.ID, cc.EventVersion, cc.State, cc.Content)
	if err != nil {
		return err
	}

	return nil
}
