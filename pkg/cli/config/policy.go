package config

import (
	"log/slog"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/opac"
	"github.com/urfave/cli/v3"
)

type Policy struct {
	path string
}

func (x *Policy) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "policy",
			Aliases:     []string{"p"},
			Usage:       "Path to policy files or directory",
			Sources:     cli.EnvVars("TRANSMITH_POLICY"),
			Destination: &x.path,
		},
	}
}

func (x Policy) LogValue() slog.Value {
	return slog.GroupValue(slog.String("path", x.path))
}

func (x Policy) New() (*opac.Client, error) {
	if x.path == "" {
		return nil, goerr.New("policy-path is not set")
	}

	return opac.New(opac.Files(x.path))
}
