# SQL export CLI [![Build Status](https://travis-ci.org/alexhokl/go-sql-export.svg?branch=master)](https://travis-ci.org/alexhokl/go-sql-export)

A command line tool to dump SQL data using T-SQL query

Usage:
  go-sql-export [command]

### Available Commands

Command | Description
---     | ---
screen  | Dump data from database and print it on screen
gsheets | Dump data from database and upload onto Google Sheets

Use "go-sql-export [command] --help" for more information about a command.

Example:
  `go-sql-export gsheets -c config.yml`

### Sample configuration

##### `config.yml`

```yml
database_type: mssql
database:
  server: example.com
  port: 1433
  name: Northwind
  username: sa
  password: pass
google_client_secret_file_path: ~/Downloads/client-secret.json
document_name: Google.DocumentExport.Example
sheets:
  - name: users
    query: "SELECT TOP 10 * FROM Users"
    columns:
      - index: 5
        data_type: date
        format: dd-MM-yyyy
```

### Installation

```sh
go install github.com/alexhokl/go-sql-export@latest
```

### Development

Please visit [Google API console](https://console.cloud.google.com/) to create
an application and enable Google Drive API. From [section Credentials in APIs
and services](https://console.cloud.google.com/apis/credentials), download
credentials file and save it. Configure the path to this file as
`google_client_secret_file_path`.

