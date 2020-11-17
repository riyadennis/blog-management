package commands

import (
	"encoding/json"
	"fmt"
	"github.com/riyadennis/blog-management/db"
	"github.com/riyadennis/blog-management/events"
	"io/ioutil"
	"net/http"
)

type CommandHandler interface {
	AggregateID() string
	Create(store db.EventStore, w http.ResponseWriter, r *http.Request)
}

type CommandArticle struct {
	Event events.Event
}

func (c *CommandArticle) AggregateID() string {
	return c.Event.AggregateID()
}

func NewCommand() *CommandArticle {
	return &CommandArticle{}
}
func (c *CommandArticle) Create(store db.EventStore, w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	a := &events.Article{}
	err = json.Unmarshal(d, a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	c.Event = events.AssignEvent(a.State, a)
	ctx := r.Context()

	err = store.Add(ctx, a.State, c.Event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	url := fmt.Sprintf("/v1/%s", c.AggregateID())
	w.Write([]byte(url))
}
