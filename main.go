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
		fmt.Println(err)
		os.Exit(1)
	}
}
