package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/riyadennis/blog-management/api/v1/article"
	"github.com/riyadennis/blog-management/pkg/api"
)

// APIv1 holds resources needed for v1 API
type APIv1 struct {
	article *article.Handler
}

// NewAPIv1 creates v1 API struct
func NewAPIv1() *APIv1 {
	return &APIv1{}
}

// ServeHTTP is the root handler for API v1
func (a *APIv1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	w.Header().Set("Content-Type", "application/json")

	resourceName, resourceParam, err := resource(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch resourceName {
	case "articles":
		article.NewHandler(resourceParam).ServeHTTP(w, r)
	default:
		http.Error(w, "invalid URL", http.StatusBadRequest)
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
