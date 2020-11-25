package eventsource

import (
	"context"
	"errors"
	"time"
)

func (c *Config) Aggregate(ctx context.Context, aggregateID string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)
	defer cancel()

	if aggregateID == "" {
		return nil, errors.New("empty aggregate id")
	}

	var content []byte
	row := c.Conn.QueryRowContext(ctx, "SELECT content as version FROM events_store WHERE aggregateID=?", aggregateID)

	err := row.Scan(&content)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, nil
	}

	return content, nil
}
