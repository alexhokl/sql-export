package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexhokl/helper/cli"
	"github.com/alexhokl/helper/database"
	"github.com/alexhokl/helper/googleapi"
	"github.com/alexhokl/sql-export/model"
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
			return runSheetExport(cmd.Context(), config, replacements)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.configFilePath, "config", "c", "", "path to configuration file")
	flags.StringArrayVarP(&opts.replacements, "replace", "r", []string{}, "replacements to SQL queries (in format key:value)")

	return cmd
}

func runSheetExport(ctx context.Context, config *model.ExportConfig, replacements map[string]string) error {
	conn, err := getDatabaseConnection(config)
	if err != nil {
		return err
	}

	dataList, err := getData(conn, config.Sheets, replacements)
	if err != nil {
		return err
	}

	errUpload := uploadDataList(ctx, dataList, config)
	if errUpload != nil {
		return errUpload
	}

	return nil
}

func uploadDataList(ctx context.Context, list []database.TableData, config *model.ExportConfig) error {
	scopes := []string{
		"https://www.googleapis.com/auth/spreadsheets",
	}
	client, errAuth := googleapi.New(
		ctx,
		config.GoogleClientSecretFilePath,
		nil,
		scopes,
	)
	if errAuth != nil {
		return errAuth
	}
	_, err := client.GetToken()
	if err != nil {
		return err
	}

	httpClient := client.NewHttpClient()
	service, errCreateService := googleapi.NewSpreadsheetService(ctx, httpClient)
	if errCreateService != nil {
		return errCreateService
	}

	document, errCreateDocument := googleapi.CreateSpreadSheet(service, config.DocumentName)
	if errCreateDocument != nil {
		return errCreateDocument
	}

	for index, data := range list {
		sheetId, errCreate := googleapi.CreateSheet(
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

		errColumns := googleapi.UpdateColumnHeaders(
			service,
			document,
			config.Sheets[index].Name,
			data.Columns,
		)
		if errColumns != nil {
			return errColumns
		}

		errRows := googleapi.UpdateRows(
			service,
			document,
			config.Sheets[index].Name,
			data.Rows,
		)
		if errRows != nil {
			return errRows
		}

		var columnConfig []googleapi.ColumnFormatConfig
		for _, c := range config.Sheets[index].Columns {
			columnConfig = append(columnConfig, c)
		}

		errFormat := googleapi.UpdateColumnStyles(
			service,
			document,
			sheetId,
			columnConfig,
		)
		if errFormat != nil {
			return errFormat
		}
	}

	fmt.Printf("spreadsheet has been created on [%s]", document.SpreadsheetUrl)
	if err := cli.OpenInBrowser(document.SpreadsheetUrl); err != nil {
		return err
	}

	return nil
}
