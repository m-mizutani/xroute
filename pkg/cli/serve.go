package cli

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/xroute/pkg/adapter"
	"github.com/m-mizutani/xroute/pkg/cli/config"
	http_server "github.com/m-mizutani/xroute/pkg/controller/http"
	"github.com/m-mizutani/xroute/pkg/usecase"
	"github.com/m-mizutani/xroute/pkg/utils/logging"
	"github.com/urfave/cli/v3"
)

func cmdServe() *cli.Command {
	var (
		addr   string
		logger config.Logger
		policy config.Policy
		slack  config.Slack
	)

	flags := joinFlags([]cli.Flag{
		&cli.StringFlag{
			Name:        "addr",
			Aliases:     []string{"a"},
			Value:       "localhost:8080",
			Usage:       "Address to listen on",
			Sources:     cli.EnvVars("XROUTE_ADDR"),
			Destination: &addr,
		},
	},
		logger.Flags(),
		policy.Flags(),
		slack.Flags(),
	)

	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Start HTTP server",
		Flags:   flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Initialize logger
			newLogger, logCloser, err := logger.New()
			if err != nil {
				return goerr.Wrap(err, "failed to create logger")
			}
			defer logCloser()
			logging.SetDefault(newLogger)

			newLogger.Info("Starting server",
				"addr", addr,
				"logger", logger,
				"policy", policy,
				"slack", slack,
			)

			var adapterOptions []adapter.Option
			if client := slack.New(); client != nil {
				adapterOptions = append(adapterOptions, adapter.WithSlack(client))
			}

			if client, err := policy.New(); err != nil {
				return goerr.Wrap(err, "failed to create policy client")
			} else {
				adapterOptions = append(adapterOptions, adapter.WithPolicy(client))
			}

			adapters := adapter.New(adapterOptions...)
			uc := usecase.New(adapters)

			s := &http.Server{
				Addr:              addr,
				ReadHeaderTimeout: 3 * time.Second,
				Handler:           http_server.New(uc),
			}

			errCh := make(chan error, 1)

			go func() {
				if err := s.ListenAndServe(); err != nil {
					errCh <- goerr.Wrap(err, "failed to listen")
				}
			}()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)

			select {
			case sig := <-sigCh:
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := s.Shutdown(ctx); err != nil {
					return goerr.Wrap(err, "failed to shutdown server", goerr.V("signal", sig))
				}

			case err := <-errCh:
				return err
			}
			return nil
		},
	}
}
