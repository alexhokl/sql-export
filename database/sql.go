package database

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// TableData struct
type TableData struct {
	Rows    [][]interface{}
	Columns []string
}

// GetConnection returns a SQL database connection
func GetConnection(connectionString string) (*sql.DB, error) {
	conn, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// GetData returns data retrieved by using query with conn
func GetData(conn *sql.DB, query string) (*TableData, error) {
	rows, errQuery := conn.Query(query)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	cols, errColumns := rows.Columns()
	if errColumns != nil {
		return nil, errColumns
	}

	columnCount := len(cols)

	var dataRows [][]interface{}
	for rows.Next() {
		vals := make([]interface{}, columnCount)
		for i := 0; i < columnCount; i++ {
			vals[i] = new(interface{})
		}
		err := rows.Scan(vals...)
		if err != nil {
			return nil, err
		}
		dataRows = append(dataRows, vals)
	}

	data := &TableData{
		Rows:    dataRows,
		Columns: cols,
	}

	return data, nil
}
