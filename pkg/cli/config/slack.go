package config

import (
	"log/slog"

	"github.com/slack-go/slack"
	"github.com/urfave/cli/v3"
)

type Slack struct {
	token string
}

func (x *Slack) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "slack-oauth-token",
			Aliases:     []string{"t"},
			Usage:       "Slack OAuth token. If empty, Slack integration is disabled",
			Value:       "",
			Sources:     cli.EnvVars("TRANSMITH_SLACK_OAUTH_TOKEN"),
			Destination: &x.token,
		},
	}
}

func (x Slack) LogValue() slog.Value {
	return slog.GroupValue(slog.Int("len(token)", len(x.token)))
}

func (x Slack) New() *slack.Client {
	if x.token == "" {
		return nil
	}
	return slack.New(x.token)
}
