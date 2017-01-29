package main

import (
	"fmt"
	"os"

	"github.com/alexhokl/go-sql-export/command"
	"github.com/alexhokl/go-sql-export/model"
)

func main() {
	config := model.ExportConfig{}

	managerCli := command.NewManagerCli(&config)
	cmd := command.NewManagerCommand(managerCli)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
