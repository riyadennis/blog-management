package commands

import (
	"github.com/riyadennis/blog-management/events"
)

type Command interface {
	AggregateID() string
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
