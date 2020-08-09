package command

import (
	"errors"

	"github.com/alexhokl/database"
	"github.com/alexhokl/go-sql-export/model"
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
				return errors.New("Configuration file is not specified")
			}
			config, errConfig := model.ParseConfig(opts.configFilePath)
			if errConfig != nil {
				return errConfig
			}
			return runScreen(config)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configFilePath, "config", "c", "", "path to configuration file")

	return cmd
}

func runScreen(config *model.ExportConfig) error {
	conn, errConn := database.GetConnection(&config.Database)
	if errConn != nil {
		return errConn
	}
	defer conn.Close()

	dataList := []database.TableData{}
	for _, s := range config.Sheets {
		data, errData := database.GetData(conn, s.Query)
		if errData != nil {
			return errData
		}
		dataList = append(dataList, *data)
	}

	database.DumpTables(dataList)

	return nil
}
