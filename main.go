package main

import (
	"fmt"
	"os"

	"github.com/alexhokl/sql-export/command"
)

func main() {
	managerCli := command.NewManagerCli()
	cmd := command.NewManagerCommand(managerCli)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
