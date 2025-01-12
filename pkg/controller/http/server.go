package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/transmith/pkg/domain/interfaces"
	"github.com/m-mizutani/transmith/pkg/domain/types"
	"github.com/m-mizutani/transmith/pkg/utils/logging"
	"github.com/m-mizutani/transmith/pkg/utils/safe"
)

type Server struct {
	router *chi.Mux
}

func New(uc interfaces.UseCases) *Server {
	r := chi.NewRouter()
	r.Use(logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		safe.Write(r.Context(), w, []byte("OK"))
	})

	r.Route("/msg", func(r chi.Router) {
		r.Post("/raw/{schema}", func(w http.ResponseWriter, r *http.Request) {
			if err := handleRawMessage(r, uc); err != nil {
				handleError(r.Context(), w, err)
				return
			}
			safe.Write(r.Context(), w, []byte("OK"))
		})

		r.Post("/pubsub/{schema}", func(w http.ResponseWriter, r *http.Request) {
			if err := handlePubSubMessage(r, uc); err != nil {
				handleError(r.Context(), w, err)
				return
			}
			safe.Write(r.Context(), w, []byte("OK"))
		})
	})

	return &Server{
		router: r,
	}
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	logging.Extract(ctx).Error("Failed to handle request", "err", err)
	code := http.StatusInternalServerError
	switch {
	case goerr.HasTag(err, types.ErrTagBadRequest):
		code = http.StatusBadRequest
	case goerr.HasTag(err, types.ErrTagUnauthorized):
		code = http.StatusUnauthorized
	}
	http.Error(w, err.Error(), code)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
