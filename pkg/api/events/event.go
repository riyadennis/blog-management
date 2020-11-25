package events

import (
	"time"
)

const (
	StatusCreated   = "created"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

type Event interface {
	Aggregate() string
	At() time.Time
}

//Model is the bounded context to hold versions of article data
// its the representation how data is store in event store
type Model struct {
	ID          string
	Version     int64
	State       string
	Content     string
	AggregateID string
	CreatedAt   time.Time
}

type EventCreated struct {
	*Model
}

type EventPublished struct {
	*Model
}

type EventArchived struct {
	*Model
}

func (a *Model) Aggregate() string {
	return a.AggregateID
}

func (a *Model) At() time.Time {
	return a.CreatedAt
}
