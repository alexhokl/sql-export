package command

import (
	"database/sql"
	"fmt"
	"sort"
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

func sortedKeys(replacements map[string]string) []string {
	keys := make([]string, 0, len(replacements))
	for k := range replacements {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] < keys[j]
	})
	return keys
}

func getData(conn *sql.DB, sheets []model.SheetConfig, replacements map[string]string) ([]database.TableData, error) {
	keys := sortedKeys(replacements)

	dataList := []database.TableData{}
	for _, s := range sheets {
		query := s.Query
		for _, k := range keys {
			query = strings.ReplaceAll(query, k, replacements[k])
		}
		data, errData := database.GetData(conn, query)
		if errData != nil {
			return nil, errData
		}
		dataList = append(dataList, *data)
	}
	return dataList, nil
}
