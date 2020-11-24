package patch

import (
	"errors"
	"github.com/riyadennis/blog-management/api/v1/article/get"
	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
	"log"
	"net/http"
	"strings"
	"time"
)

// ChangeStatus adds a new event into the store with a meta data change for the event
// in this case meta data is the status of the blog.
// I am saving aggregated record into store
// this saved aggregate can be used while querying.
func ChangeStatus( r *http.Request) error {
	path := strings.Split(r.URL.Path, "/")
	ctx := r.Context()
	refID := path[len(path)-3]

	if operation := path[len(path)-2]; operation != "status" {
		return errors.New("invalid url")
	}

	status := path[len(path)-1]

	if status != events.StatusCreated &&
		status != events.StatusPublished &&
		status != events.StatusArchived {
		return errors.New("invalid status")
	}

	store := eventsource.Get()
	version, err := store.LatestVersion(ctx, refID)
	if err != nil {
		log.Printf("failed to check version number %v", err)
		return err
	}

	content, err := get.Query(ctx, refID)
	if err != nil {
		log.Printf("failed to run query %v", err)
		return err
	}

	err = store.Apply(ctx, &events.Model{
		ID:        refID,
		Version:   version + 1,
		State:     status,
		Content:   string(content),
		Aggregate: true,
		CreatedAt: time.Now(),
	})

	return nil
}
