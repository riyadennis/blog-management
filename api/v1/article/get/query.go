package get

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

func Query(ctx context.Context, store eventsource.EventStore, refID string) ([]byte, error) {
	if store == nil {
		return nil, errors.New("empty event store config")
	}

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
