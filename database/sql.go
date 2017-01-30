package database

import "database/sql"

// GetConnection returns a SQL database connection
func GetConnection(connectionString string) (*sql.DB, error) {
	conn, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// GetData returns data retrieved by using query with conn
func GetData(conn *sql.DB, query string) ([][]interface{}, []string, error) {
	rows, errQuery := conn.Query(query)
	if errQuery != nil {
		return nil, nil, errQuery
	}
	defer rows.Close()

	cols, errColumns := rows.Columns()
	if errColumns != nil {
		return nil, nil, errColumns
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
			return nil, nil, err
		}
		dataRows = append(dataRows, vals)
	}

	return dataRows, cols, nil
}
