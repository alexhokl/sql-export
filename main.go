package main

import (
	"fmt"
	"os"

	"github.com/alexhokl/go-sql-export/command"
	"github.com/alexhokl/go-sql-export/model"
	"github.com/spf13/cobra"
)

func main() {
	config := model.ExportConfig{}

	managerCli := command.NewManagerCli(&config)
	cmd := newManagerCommand(managerCli)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newManagerCommand(cli *command.ManagerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "go-sql-export",
		Short:        "SQL data exporter",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ShowHelp(cmd, args)
		},
	}
	command.AddCommands(cmd, cli)
	return cmd
}
