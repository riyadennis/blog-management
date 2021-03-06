package post

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/riyadennis/blog-management/pkg/api/events"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"
)

// Command creates a new version of the resource for the specified ref id.
func Command(refID string, r *http.Request) error {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body %v", err)
		return err
	}

	a := &eventsource.Article{}
	err = json.Unmarshal(d, a)
	if err != nil {
		log.Printf("failed to unmarshal request body %v", err)
		return err
	}

	ctx := r.Context()

	store := eventsource.Get()

	// this is to validate and clean content
	article, err := json.Marshal(a)
	if err != nil {
		log.Printf("failed to marshal %v", err)
		return err
	}

	err = store.Apply(ctx, &events.Model{
		ID:          refID,
		State:       events.StatusCreated,
		Content:     string(article),
		AggregateID: "",
		CreatedAt:   time.Now(),
	})
	if err != nil {
		log.Printf("failed to save data to db %v", err)
		return err
	}

	return nil
}
