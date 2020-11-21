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

	eventHistory, err := store.LatestVersion(ctx, refID)
	if err != nil {
		return err
	}

	err = store.Apply(ctx, &events.Model{
		ID:        refID,
		Version:   eventHistory + 1,
		State:     events.StatusCreated,
		Content:   a,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
