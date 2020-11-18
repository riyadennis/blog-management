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
	AggregateID() string
	At() time.Time
}

//Model is the bounded context to hold versions of article data
// its the representation how data is store in event store
type Model struct {
	ID        string
	Version   int
	State     string
	Content   interface{}
	CreatedAt time.Time
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

func (a *Model) AggregateID() string {
	return a.ID
}

func (a *Model) At() time.Time {
	return a.CreatedAt
}

func AssignEvent(a *Model) Event {
	switch a.State {
	case StatusCreated:
		return &EventCreated{Model: a}
	case StatusPublished:
		return &EventPublished{Model: a}
	case StatusArchived:
		return &EventArchived{Model: a}
	default:
		return nil
	}
}
