package eventsource

import (
	"context"
	"errors"
)

// Aggregate runs through all the versions of the event and create an aggregate
// version of the event for quick fetching.
func (c *Config) Aggregate(ctx context.Context, aggregateID string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
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
