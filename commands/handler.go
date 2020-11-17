package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/riyadennis/blog-management/db"
	"github.com/riyadennis/blog-management/events"
)

// CreateArticle is the http handler which will act like a command to create a new article
func CreateArticle(store db.EventStore, w http.ResponseWriter, r *http.Request) {
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

	command := CreateCommand{
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
