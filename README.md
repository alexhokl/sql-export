# go-sql-export
A command line tool to dump SQL data using T-SQL query


#### Sample configuration

##### `config.yml`

```yml
connection_string: "server=example.com;database=Northwind;User ID=sa;Password=pass;"
document_name: Google.DocumentExport.Example
sheets:
  - name: nodes
    query: "SELECT TOP 10 * FROM Users"
```

