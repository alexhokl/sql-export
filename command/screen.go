package command

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/alexhokl/database"
	"github.com/alexhokl/go-sql-export/model"
	"github.com/spf13/cobra"
)

type screenOption struct {
	configOption
}

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

	printDataList(dataList)

	return nil
}

func printDataList(list []database.TableData) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	for _, data := range list {
		fmt.Fprintf(w, "%s\t\n", strings.Join(data.Columns, "\t"))
		for _, r := range data.Rows {
			vals := getStringValues(r)
			fmt.Fprintf(w, "%s\t\n", strings.Join(vals, "\t"))
		}
		w.Flush()
	}
	return nil
}

func getStringValues(row []interface{}) []string {
	list := []string{}
	for _, c := range row {
		list = append(list, getStringValue(c.(*interface{})))
	}
	return list
}

func getStringValue(val *interface{}) string {
	switch v := (*val).(type) {
	case nil:
		return "NULL"
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05.999")
	default:
		return fmt.Sprint(v)
	}
}
