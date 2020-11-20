package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/riyadennis/blog-management/api/v1/article"
	"github.com/riyadennis/blog-management/pkg/api"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

type APIv1 struct {
	eventStore eventsource.EventStore
	article    *article.Handler
}

func NewAPIv1(e eventsource.EventStore) *APIv1 {
	return &APIv1{eventStore: e}
}

func (a *APIv1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	w.Header().Set("Content-Type", "application/json")

	resourceName, resourceParam, err := resource(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch resourceName {
	case "article":
		article.NewHandler(a.eventStore, resourceParam).ServeHTTP(w, r)
	}
}

func resource(url string) (string, string, error) {
	head, tail := api.ParsePath(url)
	if head != "api" {
		return "", "", errors.New("invalid URL")
	}

	version, resource := api.ParsePath(tail)
	if version != "v1" {
		return "", "", errors.New("unsupported api version")
	}

	resourceName, resourceParam := api.ParsePath(resource)

	resourceParam = strings.Replace(resourceParam, "/", "", -1)

	return resourceName, resourceParam, nil
}
