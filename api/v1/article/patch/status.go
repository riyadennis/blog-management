package patch

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/riyadennis/blog-management/api/v1/article/get"
	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

// ChangeStatus adds a new event into the store with a meta data change for the event
// in this case meta data is the status of the blog.
// I am saving aggregated record into store
// this saved aggregate can be used while querying.
func ChangeStatus(r *http.Request) (string, error) {
	path := strings.Split(r.URL.Path, "/")
	ctx := r.Context()
	refID := path[len(path)-3]

	if operation := path[len(path)-2]; operation != "status" {
		return "", errors.New("invalid url")
	}

	status := path[len(path)-1]

	if status != events.StatusCreated &&
		status != events.StatusPublished &&
		status != events.StatusArchived {
		return "", errors.New("invalid status")
	}

	store := eventsource.Get()

	content, err := get.Query(r, refID)
	if err != nil {
		log.Printf("failed to run query %v", err)
		return "", err
	}

	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	aggregateID := u.String()

	err = store.Apply(ctx, &events.Model{
		ID:          refID,
		State:       status,
		Content:     string(content),
		AggregateID: aggregateID,
		CreatedAt:   time.Now(),
	})

	return aggregateID, nil
}
