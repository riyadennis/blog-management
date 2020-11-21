package get

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

type Query struct {
	EventStore eventsource.EventStore
}

func Aggregate(ctx context.Context, store eventsource.EventStore, refID string) ([]byte, error) {
	if store == nil {
		return nil, errors.New("empty event store config")
	}

	latest, err := store.LatestVersion(ctx, refID)
	if err != nil {
		return nil, err
	}

	refIDEvents, err := store.Load(ctx, refID)
	if err != nil {
		return nil, err
	}

	var eventLatest *events.Model
	for _, e := range refIDEvents {
		m, ok := e.(*events.Model)
		if !ok {
			return nil, errors.New("invalid event found in history")
		}
		//TODO append fields also to create an aggregate
		if m.Version == latest {
			eventLatest = m
		}
	}

	if eventLatest.Content == nil {
		return nil, errors.New("empty content")
	}

	data, err := json.Marshal(eventLatest.Content)
	if err != nil {
		return nil, err
	}

	return data, nil
}
