package oauth

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Authorize(w http.ResponseWriter,r *http.Request) {
	connectionID := chi.URLParam(r, "connectionID")

	redirectURL, err := h.service.BeginAuthorization(
		r.Context(),
		connectionID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) Callback(w http.ResponseWriter,r *http.Request) {

	err := h.service.HandleCallback(
		r.Context(),
		r.URL.Query(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	

	//http.Redirect(w, r, "/connections", http.StatusFound)
}

