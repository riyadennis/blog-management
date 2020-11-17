package handlers

import (
	"encoding/json"
	"github.com/riyadennis/blog-management/commands"
	"github.com/riyadennis/blog-management/events"
	"io/ioutil"
	"net/http"
	"time"
)

// CreateEvent is the http handler which will call command the to
// create a new article
// I have seen implementations where single http handler handles
// different commands by reading domain-model from the header
// TODO check with David whether what is the right approach
// command := r.Header.Get("domain-model")
// if command != "CreateEventCommand"{
//	call create event handler
// }
func CreateEvent(store commands.EventStore, w http.ResponseWriter, r *http.Request) {
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

	if a.State == "" {
		a.State = events.StatusCreated
	}

	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	command := commands.CreateCommand{
		ArticleCreated: events.ArticleCreated{Article: a},
	}

	ctx := r.Context()
	es, err := command.Apply(store, ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(es)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
