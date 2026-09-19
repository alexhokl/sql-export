package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      ExportConfig
		wantErr     bool
		errContains string
	}{
		{
			name:   "no columns",
			config: ExportConfig{Sheets: []SheetConfig{{Name: "users", Query: "SELECT 1"}}},
		},
		{
			name: "valid data types",
			config: ExportConfig{Sheets: []SheetConfig{
				{
					Name:  "users",
					Query: "SELECT 1",
					Columns: []ColumnConfig{
						{Index: 5, DataType: "date", Format: "dd-MM-yyyy"},
						{Index: 6, DataType: "money"},
					},
				},
			}},
		},
		{
			name: "unsupported data type",
			config: ExportConfig{Sheets: []SheetConfig{
				{
					Name:  "users",
					Query: "SELECT 1",
					Columns: []ColumnConfig{
						{Index: 5, DataType: "number"},
					},
				},
			}},
			wantErr:     true,
			errContains: "unsupported data_type [number] on sheet [users] column index [5]",
		},
		{
			name: "error mentions failing sheet and column",
			config: ExportConfig{Sheets: []SheetConfig{
				{
					Name:    "valid",
					Query:   "SELECT 1",
					Columns: []ColumnConfig{{Index: 1, DataType: "date"}},
				},
				{
					Name:  "orders",
					Query: "SELECT 1",
					Columns: []ColumnConfig{
						{Index: 2, DataType: "money"},
						{Index: 3, DataType: "percent"},
					},
				},
			}},
			wantErr:     true,
			errContains: "unsupported data_type [percent] on sheet [orders] column index [3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		path := writeTempConfig(t, `
database_type: mssql
database:
  server: example.com
  port: 1433
  name: Northwind
  username: sa
  password: pass
document_name: Google.DocumentExport.Example
sheets:
  - name: users
    query: "SELECT TOP 10 * FROM Users"
    columns:
      - index: 5
        data_type: date
        format: dd-MM-yyyy
`)
		config, err := ParseConfig(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if config.DatabaseType != "mssql" {
			t.Errorf("expected database_type mssql, got %q", config.DatabaseType)
		}
		if len(config.Sheets) != 1 {
			t.Fatalf("expected 1 sheet, got %d", len(config.Sheets))
		}
		if config.Sheets[0].Columns[0].Format != "dd-MM-yyyy" {
			t.Errorf("expected format dd-MM-yyyy, got %q", config.Sheets[0].Columns[0].Format)
		}
	})

	t.Run("invalid column type rejected", func(t *testing.T) {
		path := writeTempConfig(t, `
database_type: mssql
sheets:
  - name: users
    query: "SELECT 1"
    columns:
      - index: 0
        data_type: number
`)
		_, err := ParseConfig(path)
		if err == nil {
			t.Fatal("expected error for unsupported data_type, got nil")
		}
		if !strings.Contains(err.Error(), "unsupported data_type [number]") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := ParseConfig(filepath.Join(t.TempDir(), "missing.yml"))
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("malformed yaml", func(t *testing.T) {
		path := writeTempConfig(t, ":\n\tinvalid yaml content: [")
		_, err := ParseConfig(path)
		if err == nil {
			t.Fatal("expected error for malformed yaml, got nil")
		}
	})

	t.Run("duplicate keys rejected by yaml v3", func(t *testing.T) {
		path := writeTempConfig(t, `
database_type: mssql
database_type: postgres
`)
		_, err := ParseConfig(path)
		if err == nil {
			t.Fatal("expected error for duplicate keys, got nil")
		}
	})
}
