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
