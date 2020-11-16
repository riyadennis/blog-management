package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/riyadennis/blog-management/commands"
	"github.com/riyadennis/blog-management/db"
	"github.com/riyadennis/blog-management/events"
	"io/ioutil"
	"net/http"
	"time"
)

type CommandHandler interface {
	Apply(ctx context.Context, command commands.Command) ([]events.Event, error)
}

type Handler struct {
	Name string
}

func CreateArticle() *Handler {
	return &Handler{
		Name: events.StatusCreated,
	}
}

func (h *Handler) Apply(ctx context.Context, e commands.Command) ([]events.Event, error) {
	var ev []events.Event
	switch h.Name {
	case events.StatusCreated:
		ctx, cancel := context.WithTimeout(ctx, time.Second*100)
		defer cancel()
		v := e.(commands.CreateCommand)
		err := db.Store(ctx, &v)
		if err != nil {
			return nil, err
		}
		ev = append(ev, v)
		return ev, nil
	}

	return nil, fmt.Errorf("invalid command")
}

// CreateEvent is the http handler which will call command handler to
// create a new article
// I have seen implementations where single http handler handles
// different commands by reading domain-model from the header
// TODO check with David whether what is the right approach
// command := r.Header.Get("domain-model")
// if command != "CreateEventCommand"{
//	call create event handler
// }
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	a := &events.Article{}
	err = json.Unmarshal(d, a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	a.State = events.StatusCreated
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	command := commands.CreateCommand{
		ArticleCreated: events.ArticleCreated{Article: a},
	}

	ctx := r.Context()
	es, err := CreateArticle().Apply(ctx, command)
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
