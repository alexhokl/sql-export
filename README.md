# go-sql-export
A command line tool to dump SQL data using T-SQL query

Usage:
  go-sql-export [command]

Available Commands:

Command | Description
--- | ---
screen    | Dump data from database and print it on screen
gsheets     | Dump data from database and upload onto Google Sheets

Use "go-sql-export [command] --help" for more information about a command.


#### Sample configuration

##### `config.yml`

```yml
connection_string: "server=example.com;database=Northwind;User ID=sa;Password=pass;"
document_name: Google.DocumentExport.Example
sheets:
  - name: users
    query: "SELECT TOP 10 * FROM Users"
    columns:
      - index: 5
        data_type: date
        format: dd-MM-yyyy
```

