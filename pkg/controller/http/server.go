package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/xroute/pkg/domain/interfaces"
	"github.com/m-mizutani/xroute/pkg/domain/types"
	"github.com/m-mizutani/xroute/pkg/utils/logging"
	"github.com/m-mizutani/xroute/pkg/utils/safe"
)

type Server struct {
	router              *chi.Mux
	githubWebhookSecret string
}

type Option func(*Server)

func WithGitHubWebhookSecret(secret string) Option {
	return func(s *Server) {
		s.githubWebhookSecret = secret
	}
}

func New(uc interfaces.UseCases, options ...Option) *Server {
	r := chi.NewRouter()
	server := &Server{
		router: r,
	}

	for _, opt := range options {
		opt(server)
	}

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

		r.Route("/github", func(r chi.Router) {
			r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
				if err := handleGitHubWebhook(r, uc, server.githubWebhookSecret); err != nil {
					handleError(r.Context(), w, err)
					return
				}
				safe.Write(r.Context(), w, []byte("OK"))
			})

			r.Post("/actions", func(w http.ResponseWriter, r *http.Request) {
				if err := handleGitHubActions(r, uc); err != nil {
					handleError(r.Context(), w, err)
					return
				}
				safe.Write(r.Context(), w, []byte("OK"))
			})
		})
	})

	return server
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
