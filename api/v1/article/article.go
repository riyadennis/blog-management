package article

import (
	"github.com/riyadennis/blog-management/api/v1/article/create"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
	"net/http"
	"path"
)

func Dispatch(url string, store eventsource.EventStore, r *http.Request) error {
	_, refID := path.Split(url)
	err := create.ArticleEvent(store, refID, r)
	if err != nil {
		return err
	}

	return nil
}
