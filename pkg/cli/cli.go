package cli

import (
	"context"

	"github.com/m-mizutani/xroute/pkg/utils/logging"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, args []string) error {
	app := cli.Command{
		Name:  "xroute",
		Usage: "Manipulate and transmit Webhook messages by Rego policies",
		Commands: []*cli.Command{
			cmdServe(),
		},
	}

	if err := app.Run(ctx, args); err != nil {
		logging.Default().Error("Failed to run command", "error", err)
		return err
	}

	return nil
}

func joinFlags(flags ...[]cli.Flag) []cli.Flag {
	var ret []cli.Flag
	for _, f := range flags {
		ret = append(ret, f...)
	}
	return ret
}
