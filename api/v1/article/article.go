package article

import (
	"github.com/riyadennis/blog-management/api/v1/article/get"
	"github.com/riyadennis/blog-management/api/v1/article/patch"
	"github.com/riyadennis/blog-management/api/v1/article/post"
	"net/http"
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
			w.WriteHeader(http.StatusBadRequest)
			resp = []byte(err.Error())
			return
		}
	case http.MethodPatch:
		id, err := patch.ChangeStatus(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp = []byte(err.Error())
			return
		}
		w.Header().Set("Etag", id)
	case http.MethodGet:
		event, err := get.Query(r, h.resourceID)
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
