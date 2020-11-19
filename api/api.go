package api

import (
	"github.com/riyadennis/blog-management/api/v1/article"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
	"net/http"
	"path"
	"strings"
)

type APIv1 struct {
	eventStore eventsource.EventStore
}

func NewAPIv1(e eventsource.EventStore) *APIv1 {
	return &APIv1{eventStore: e}
}

func (a *APIv1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	w.Header().Set("Content-Type", "application/json")
	matched, err := path.Match("/api/v1/*/*", p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid URL"))
		return
	}

	if !strings.Contains(p, "article") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid URL"))
		return
	}

	err = article.Dispatch(p, a.eventStore, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
