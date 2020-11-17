package commands

import (
	"context"
	"github.com/riyadennis/blog-management/events"
	"time"
)

type Command interface {
	AggregateID() string
	Apply(store EventStore, ctx context.Context) ([]events.Event, error)
}

type CreateCommand struct {
	events.ArticleCreated
}

func (c CreateCommand) AggregateID() string {
	return c.ID
}

func (c CreateCommand) Status() string {
	return c.State
}

func (c CreateCommand) Apply(store EventStore, ctx context.Context) ([]events.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	err := store.Add(ctx, c)
	if err != nil {
		return nil, err
	}

	return nil, nil

}
