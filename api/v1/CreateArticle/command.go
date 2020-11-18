package CreateArticle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

type Command struct {
	EventStore eventsource.EventStore
	Event      events.Event
}

func (c *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	a := &events.Model{}
	err = json.Unmarshal(d, a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jc, err := json.Marshal(a.Content)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ctx := r.Context()

	eventHistory, err := c.EventStore.Load(ctx, a.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	latestVersion := eventVersion(eventHistory)

	a.Version = latestVersion + 1
	c.SetEvent(events.AssignEvent(a))

	err = c.EventStore.Apply(ctx, a, string(jc))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	url := fmt.Sprintf("/v1/%s", c.AggregateID())
	w.Write([]byte(url))
}

func NewCommand(e eventsource.EventStore) *Command {
	return &Command{EventStore: e}
}

func (c *Command) AggregateID() string {
	return c.Event.AggregateID()
}

func (c *Command) SetEvent(e events.Event) {
	c.Event = e
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
	return latestVersion
}
