package commands

import (
	"encoding/json"
	"fmt"
	"github.com/riyadennis/blog-management/events"
	"github.com/riyadennis/blog-management/eventsource"
	"io/ioutil"
	"net/http"
)

type CommandHandler interface {
	SetEvent(e events.Event)
	GetEvent() events.Event
	AggregateID() string
	CreateArticle(store eventsource.EventStore, w http.ResponseWriter, r *http.Request)
}

type CommandArticle struct {
	Event events.Event
}

func (c *CommandArticle) SetEvent(e events.Event) {
	c.Event = e
}

func (c *CommandArticle) GetEvent() events.Event {
	if c == nil {
		return nil
	}

	return c.Event
}

func (c *CommandArticle) AggregateID() string {
	return c.Event.AggregateID()
}

func NewCommand() *CommandArticle {
	return &CommandArticle{}
}

func (c *CommandArticle) CreateArticle(eventStore eventsource.EventStore, w http.ResponseWriter, r *http.Request) {
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
	ctx := r.Context()
	eventHistory, err := eventStore.Load(ctx, a.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	latestVersion := eventVersion(eventHistory)

	a.Version = latestVersion + 1

	c.SetEvent(events.AssignEvent(a))

	err = eventStore.Apply(ctx, a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	url := fmt.Sprintf("/v1/%s", c.AggregateID())
	w.Write([]byte(url))
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
