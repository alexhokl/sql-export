package command

import (
	"errors"

	"github.com/alexhokl/helper/database"
	"github.com/alexhokl/sql-export/model"
	"github.com/spf13/cobra"
)

type screenOption struct {
	configOption
}

// NewScreenCommand returns a command
func NewScreenCommand(cli *ManagerCli) *cobra.Command {
	opts := screenOption{}

	cmd := &cobra.Command{
		Use:   "screen",
		Short: "Export data on-screen",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				cli.ShowHelp(cmd, args)
				return nil
			}
			if opts.configFilePath == "" {
				return errors.New("configuration file is not specified")
			}
			config, errConfig := model.ParseConfig(opts.configFilePath)
			if errConfig != nil {
				return errConfig
			}
			replacements, err := getReplacementMap(opts.replacements)
			if err != nil {
				return err
			}
			return runScreen(config, replacements)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configFilePath, "config", "c", "", "path to configuration file")
	flags.StringArrayVarP(&opts.replacements, "replace", "r", []string{}, "replacements to SQL queries (in format key:value)")

	return cmd
}

func runScreen(config *model.ExportConfig, replacements map[string]string) error {
	conn, errConn := getDatabaseConnection(config)
	if errConn != nil {
		return errConn
	}
	// defer conn.Close()

	dataList, err := getData(conn, config.Sheets, replacements)
	if err != nil {
		return err
	}

	database.DumpTables(dataList)

	return nil
}
