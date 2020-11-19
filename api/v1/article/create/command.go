package create

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

func ArticleEvent(store eventsource.EventStore, refID string, r *http.Request) error {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	a := &eventsource.Article{}
	err = json.Unmarshal(d, a)
	if err != nil {
		return err
	}

	ctx := r.Context()

	eventHistory, err := store.Load(ctx, refID)
	if err != nil {
		return err
	}

	err = store.Apply(ctx, &events.Model{
		ID: refID,
		// TODO find a way to do auto increment in db
		// this can become a resource consuming process
		Version:   eventVersion(eventHistory),
		State:     events.StatusCreated,
		Content:   a,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func eventVersion(eventHistory []events.Event) int {
	if eventHistory == nil {
		return 1
	}
	var latestVersion int
	for _, e := range eventHistory {
		if e == nil {
			continue
		}
		m := e.(*events.Model)
		if m.Version > latestVersion {
			latestVersion = m.Version
		}
	}
	return latestVersion + 1
}
