package patch

import (
	"errors"
	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
	"log"
	"net/http"
	"strings"
	"time"
)

func ChangeStatus(store eventsource.EventStore, r *http.Request) error {
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

	version, err := store.LatestVersion(ctx, refID)
	if err != nil {
		log.Printf("failed to check version number %v", err)
		return err
	}

	err = store.Apply(ctx, &events.Model{
		ID:        refID,
		Version:   version + 1,
		State:     status,
		Content:   nil,
		CreatedAt: time.Now(),
	})

	return nil
}