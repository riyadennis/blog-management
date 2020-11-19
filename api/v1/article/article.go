package article

import (
	"github.com/riyadennis/blog-management/api/v1/article/create"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
	"net/http"
	"path"
)

type Article struct {
	eventStore eventsource.EventStore
}

func NewArticle(e eventsource.EventStore) *Article {
	return &Article{eventStore: e}
}

func (a *Article) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	w.Header().Set("Content-Type", "application/json")
	matched, err := path.Match("/api/v1/article/*", p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid method"))
		return
	}

	if matched {
		_, refID := path.Split(p)
		err = create.ArticleEvent(a.eventStore, refID, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
			return
		}
	}

}
