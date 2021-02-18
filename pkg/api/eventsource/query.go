package eventsource

import (
	"context"
	"errors"
	"time"

	"github.com/riyadennis/blog-management/pkg/api/events"
)

// Load fetches all the events from db as per the aggregate ID
func (c *Config) Load(ctx context.Context, aggregateID string) ([]events.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	if aggregateID == "" {
		return nil, errors.New("empty aggregate id")
	}

	rows, err := c.Conn.QueryContext(
		ctx,
		"SELECT version,state,content,created_at FROM events_store WHERE resourceID=? ORDER BY created_at DESC",
		aggregateID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var models []events.Event

	for rows.Next() {
		var state, data, createdAt string
		var version int64

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
