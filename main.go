package main

import (
	"context"
	"os"

	"github.com/m-mizutani/xroute/pkg/cli"
)

func main() {
	ctx := context.Background()
	if err := cli.Run(ctx, os.Args); err != nil {
		os.Exit(1)
	}
}
