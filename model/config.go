package model

import (
	"io/ioutil"
	"path/filepath"

	"github.com/alexhokl/helper/database"
	yaml "gopkg.in/yaml.v2"
)

// ExportConfig struct
type ExportConfig struct {
	DatabaseType               string          `yaml:"database_type"`
	Database                   database.Config `yaml:"database"`
	DocumentName               string          `yaml:"document_name"`
	Sheets                     []SheetConfig   `yaml:"sheets"`
	GoogleClientSecretFilePath string          `yaml:"google_client_secret_file_path"`
}

// SheetConfig struct
type SheetConfig struct {
	Name    string         `yaml:"name"`
	Query   string         `yaml:"query"`
	Columns []ColumnConfig `yaml:"columns,omitempty"`
}

// ColumnConfig struct
type ColumnConfig struct {
	Index    int    `yaml:"index"`
	DataType string `yaml:"data_type"`
	Format   string `yaml:"format"`
}

func (c ColumnConfig) ColumnIndex() int {
	return c.Index
}

func (c ColumnConfig) ColumnDataType() string {
	return c.DataType
}

func (c ColumnConfig) ColumnFormat() string {
	return c.Format
}

// ParseConfig returns configuration from the specified file path
func ParseConfig(filePath string) (*ExportConfig, error) {
	filename, errPath := filepath.Abs(filePath)
	if errPath != nil {
		return nil, errPath
	}
	yamlFile, errIo := ioutil.ReadFile(filename)
	if errIo != nil {
		return nil, errIo
	}
	config := ExportConfig{}
	err := yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
