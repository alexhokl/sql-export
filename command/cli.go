package command

import (
	"github.com/alexhokl/go-sql-export/model"
	"github.com/spf13/cobra"
)

// ManagerCli struct
type ManagerCli struct {
	config *model.ExportConfig
}

// NewManagerCli creates a new manager cli instance
func NewManagerCli(config *model.ExportConfig) *ManagerCli {
	cli := ManagerCli{
		config: config,
	}
	return &cli
}

// ShowHelp shows the command help
func (cli *ManagerCli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.HelpFunc()(cmd, args)
	return nil
}

type option struct {
	configFilePath string
}

// NewManagerCommand returns the main command of this exporter
func NewManagerCommand(cli *ManagerCli) *cobra.Command {
	opts := option{}

	cmd := &cobra.Command{
		Use:          "go-sql-export",
		Short:        "SQL data exporter",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ShowHelp(cmd, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configFilePath, "config", "c", "", "path to configuration file")

	AddCommands(cmd, cli, opts)
	return cmd
}
