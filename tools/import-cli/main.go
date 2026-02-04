package main

import (
	"os"

	"github.com/criblio/terraform-provider-criblio/tools/import-cli/cmd"
)

func main() {
	root := cmd.NewRootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
