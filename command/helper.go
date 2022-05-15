package command

import (
	"database/sql"
	"fmt"

	"github.com/alexhokl/go-sql-export/model"
	"github.com/alexhokl/helper/database"
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

func getData(conn *sql.DB, sheets []model.SheetConfig) ([]database.TableData, error) {
	dataList := []database.TableData{}
	for _, s := range sheets {
		data, errData := database.GetData(conn, s.Query)
		if errData != nil {
			return nil, errData
		}
		dataList = append(dataList, *data)
	}
	return dataList, nil
}
