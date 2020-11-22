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

	latest, err := store.LatestVersion(ctx, refID)
	if err != nil {
		return nil, err
	}

	refIDEvents, err := store.Load(ctx, refID)
	if err != nil {
		return nil, err
	}

	eventLatest, err := aggregate(refIDEvents, latest)
	if err != nil {
		return nil, err
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

func aggregate(ev []events.Event, v int64) (*events.Model, error) {
	eventLatest := &events.Model{}

	eventLatest.Version = v

	for i, e := range ev {
		m, ok := e.(*events.Model)
		if !ok {
			return nil, errors.New("invalid event found in history")
		}
		if m.Version == v {
			eventLatest = m
		}

		if m.Content == "null" {
			m1, ok := ev[i+1].(*events.Model)
			if !ok {
				return nil, errors.New("invalid event found in history")
			}

			eventLatest.Content = m1.Content
		}
	}

	return eventLatest, nil
}

func recursive(e *events.Model, count int, es []events.Event) error {
	m, ok := es[count+1].(*events.Model)
	if !ok {
		return errors.New("invalid event found in history")
	}

	e.Content = m.Content

	return nil
}
