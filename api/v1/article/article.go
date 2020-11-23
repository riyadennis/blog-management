package article

import (
	"github.com/riyadennis/blog-management/api/v1/article/get"
	"github.com/riyadennis/blog-management/api/v1/article/patch"
	"github.com/riyadennis/blog-management/api/v1/article/post"
	"github.com/riyadennis/blog-management/pkg/api/eventsource"

	"net/http"
)

type Handler struct {
	store      eventsource.EventStore
	resourceID string
}

func NewHandler(store eventsource.EventStore, resourceID string) *Handler {
	return &Handler{
		store:      store,
		resourceID: resourceID,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	switch r.Method {
	case http.MethodPost:
		err := post.Command(h.store, h.resourceID, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp = []byte(err.Error())
			return
		}
	case http.MethodPatch:
		err := patch.ChangeStatus(h.store, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp = []byte(err.Error())
			return
		}
	case http.MethodGet:
		event, err := get.Query(r.Context(), h.store, h.resourceID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if event == nil {
			w.WriteHeader(http.StatusBadRequest)
			resp = []byte("no event found for provided resource ID")
			w.Write(resp)
			return
		}
		resp = event
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
