package command

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alexhokl/helper/database"
	"github.com/alexhokl/sql-export/model"
)

func getDatabaseConnection(config *model.ExportConfig) (*sql.DB, error) {
	switch config.DatabaseType {
	case "mssql":
		return database.GetConnection(&config.Database)
	case "postgres":
		c := &database.PostgresConfig{
			Config: config.Database,
			UseSSL: true,
		}
		return database.GetPostgresConnection(c)
	default:
		return nil, fmt.Errorf("un-supported database type [%s]", config.DatabaseType)
	}
}

func getData(conn *sql.DB, sheets []model.SheetConfig, replacements map[string]string) ([]database.TableData, error) {
	dataList := []database.TableData{}
	for _, s := range sheets {
		query := s.Query
		for k, v := range replacements {
			query = strings.ReplaceAll(query, k, v)
		}
		data, errData := database.GetData(conn, query)
		if errData != nil {
			return nil, errData
		}
		dataList = append(dataList, *data)
	}
	return dataList, nil
}
