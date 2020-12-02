package article

import (
	"net/http"

	"github.com/riyadennis/blog-management/api/v1/article/get"
	"github.com/riyadennis/blog-management/api/v1/article/patch"
	"github.com/riyadennis/blog-management/api/v1/article/post"
)

type Handler struct {
	resourceID string
}

func NewHandler(resourceID string) *Handler {
	return &Handler{
		resourceID: resourceID,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	switch r.Method {
	case http.MethodPost:
		err := post.Command(h.resourceID, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case http.MethodPatch:
		id, err := patch.ChangeStatus(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Etag", id)
	case http.MethodGet:
		event, err := get.Query(r, h.resourceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if event == nil {
			http.Error(w, "no event found for provided resource ID", http.StatusBadRequest)
			return
		}
		resp = event
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
