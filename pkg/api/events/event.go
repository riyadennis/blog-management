package events

import (
	"time"
)

const (
	// StatusCreated is the status for the event in db when an article is created
	StatusCreated = "created"

	// StatusPublished is the status for the event in db when an article is published
	// change to this status happens through a call from a separate end point
	StatusPublished = "published"

	// StatusArchived is the status for the event in db when an article is archived
	// change to this status happens through a call from a separate end point
	StatusArchived = "archived"
)

// Event holds contract to interact with event source
type Event interface {
	Aggregate() string
	At() time.Time
}

// Model is the bounded context to hold versions of article data
// its the representation how data is store in event store
type Model struct {
	ID          string
	Version     int64
	State       string
	Content     string
	AggregateID string
	CreatedAt   time.Time
}

// Aggregate will aggregate id to be used to fetch the resource
func (a *Model) Aggregate() string {
	return a.AggregateID
}

// At is the getter for created record
func (a *Model) At() time.Time {
	return a.CreatedAt
}
