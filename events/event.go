package events

import "time"

const (
	StatusCreated   = "created"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

type Event interface {
	AggregateID() string
	Version() int
	At() time.Time
	Data() string
}

//Article is the bounded context to hold versions of article data
type ArticleEvent struct {
	ID           string
	EventVersion int
	State        string
	Article      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ArticleCreated struct {
	*ArticleEvent
}

type ArticlePublished struct {
	*ArticleEvent
}

type ArticleArchived struct {
	*ArticleEvent
}

func (a *ArticleEvent) AggregateID() string {
	return a.ID
}

func (a *ArticleEvent) Version() int {
	return a.EventVersion
}

func (a *ArticleEvent) At() time.Time {
	return a.CreatedAt
}

func (a *ArticleEvent) Data() string {
	return a.Article
}
