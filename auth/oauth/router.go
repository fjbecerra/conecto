package oauth

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(oauth *Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Route("/oauth", func(sub chi.Router) {

		sub.Get(
			"/{connectionID}/authorize",
			oauth.Authorize,
		)

		sub.Get(
			"/callback",
			oauth.Callback,
		)
	})

	return r
}