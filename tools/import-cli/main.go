package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	root := cmd.NewRootCommand()
	if err := root.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
