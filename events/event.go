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
type Article struct {
	ID           string
	EventVersion int
	State        string
	Content      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ArticleCreated struct {
	*Article
}

type ArticlePublished struct {
	*Article
}

type ArticleArchived struct {
	*Article
}

func (a *Article) AggregateID() string {
	return a.ID
}

func (a *Article) Version() int {
	return a.EventVersion
}

func (a *Article) At() time.Time {
	return a.CreatedAt
}

func (a *Article) Data() string {
	return a.Content
}
