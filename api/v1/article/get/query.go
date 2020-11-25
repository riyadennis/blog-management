package get

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

// Query will fetch current projection of the resource based on events in the store.
func Query(r *http.Request, refID string) ([]byte, error) {
	ctx := r.Context()
	store := eventsource.Get()
	if r.Header.Get("If-None-Match") != "" {
		return store.Aggregate(ctx, r.Header.Get("If-None-Match"))
	}

	return Article(ctx, refID)
}

func Article(ctx context.Context, refID string) ([]byte, error) {
	store := eventsource.Get()
	refIDEvents, err := store.Load(ctx, refID)
	if err != nil {
		return nil, err
	}

	articles, err := allArticles(refIDEvents)
	if err != nil {
		return nil, err
	}

	if articles == nil {
		return nil, errors.New("no content in the history")
	}

	article, err := aggregate(articles)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(article)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// aggregate is a simple function which goes through all ordered articles
// and add latest fields to one article object.
// this can be done using aggregateID or flag as well which will be more performant.
// this is a simple approach.
func aggregate(articles []*eventsource.Article) (*eventsource.Article, error) {
	article := &eventsource.Article{}

	for _, a := range articles {
		if a == nil {
			continue
		}

		if a.Introduction != "" && article.Introduction == "" {
			article.Introduction = a.Introduction
		}
		if a.Heading != "" && article.Heading == "" {
			article.Heading = a.Heading
		}
		if a.Body != "" && article.Body == "" {
			article.Body = a.Body
		}
		if a.Author != "" && article.Author == "" {
			article.Author = a.Author
		}
	}

	return article, nil
}

// allArticles fetches non empty content from all the events in the store
// articles are fetched in the order of creation.
func allArticles(ev []events.Event) ([]*eventsource.Article, error) {
	ar := make([]*eventsource.Article, len(ev))

	for i, e := range ev {
		m, ok := e.(*events.Model)
		if !ok {
			return nil, errors.New("invalid event found in history")
		}

		if m.Content == "null" {
			continue
		}
		ar[i] = &eventsource.Article{}
		str := m.Content

		err := json.Unmarshal([]byte(str), ar[i])
		if err != nil {
			return nil, err
		}
	}

	return ar, nil
}
