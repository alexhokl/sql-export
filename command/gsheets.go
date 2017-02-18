package command

import (
	"errors"

	"github.com/alexhokl/go-sql-export/database"
	"github.com/alexhokl/go-sql-export/google"
	"github.com/alexhokl/go-sql-export/model"
	"github.com/spf13/cobra"
)

type gsheetsOption struct {
	configOption
}

func NewGSheetsCommand(cli *ManagerCli) *cobra.Command {
	opts := gsheetsOption{}

	cmd := &cobra.Command{
		Use:   "gsheets",
		Short: "Export data onto a Google Sheets",
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
			return runSheetExport(config)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configFilePath, "config", "c", "", "path to configuration file")

	return cmd
}

func runSheetExport(config *model.ExportConfig) error {
	conn, errConn := database.GetConnection(config.ConnectionString)
	if errConn != nil {
		return errConn
	}

	dataList := []database.TableData{}
	for _, s := range config.Sheets {
		data, errData := database.GetData(conn, s.Query)
		if errData != nil {
			return errData
		}
		dataList = append(dataList, *data)
	}

	errUpload := uploadDataList(dataList, config)
	if errUpload != nil {
		return errUpload
	}

	return nil
}

func uploadDataList(list []database.TableData, config *model.ExportConfig) error {
	httpClient, errAuth := google.NewHttpClient("client_secret.json")
	if errAuth != nil {
		return errAuth
	}

	service, errCreateService := google.NewSpreadsheetService(httpClient)
	if errCreateService != nil {
		return errCreateService
	}

	document, errCreateDocument := google.CreateSpreadSheet(service, config.DocumentName)
	if errCreateDocument != nil {
		return errCreateDocument
	}

	for index, data := range list {
		sheetId, errCreate := google.CreateSheet(
			service,
			document,
			index,
			config.Sheets[index].Name,
			data.Rows,
			data.Columns,
		)
		if errCreate != nil {
			return errCreate
		}

		errColumns := google.UpdateColumnHeaders(
			service,
			document,
			config.Sheets[index].Name,
			data.Columns,
		)
		if errColumns != nil {
			return errColumns
		}

		errRows := google.UpdateRows(
			service,
			document,
			config.Sheets[index].Name,
			data.Rows,
		)
		if errRows != nil {
			return errRows
		}

		errFormat := google.UpdateColumnStyles(
			service,
			document,
			sheetId,
			config.Sheets[index].Columns,
		)
		if errFormat != nil {
			return errFormat
		}
	}

	return nil
}
