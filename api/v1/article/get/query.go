package get

import (
	"net/http"
	"strings"

	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

type Query struct {
	EventStore eventsource.EventStore
}

func (c *Query) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	refID := strings.Split("/", r.URL.Path)
	w.Write([]byte(refID[len(refID)-1]))
}
