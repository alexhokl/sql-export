package command

import (
	"strings"
	"testing"
)

func TestGetReplacementMap(t *testing.T) {
	tests := []struct {
		name        string
		replacement []string
		expected    map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:        "simple key value",
			replacement: []string{"key:value"},
			expected:    map[string]string{"key": "value"},
		},
		{
			name:        "value containing colons",
			replacement: []string{"duration:PT30S"},
			expected:    map[string]string{"duration": "PT30S"},
		},
		{
			name:        "value is a URL",
			replacement: []string{"endpoint:https://example.com"},
			expected:    map[string]string{"endpoint": "https://example.com"},
		},
		{
			name:        "empty entries skipped",
			replacement: []string{"", "key:value"},
			expected:    map[string]string{"key": "value"},
		},
		{
			name:        "empty value allowed",
			replacement: []string{"key:"},
			expected:    map[string]string{"key": ""},
		},
		{
			name:        "missing separator",
			replacement: []string{"novalue"},
			wantErr:     true,
			errContains: `invalid replacement "novalue"`,
		},
		{
			name:        "duplicated key",
			replacement: []string{"key:one", "key:two"},
			wantErr:     true,
			errContains: `duplicated replacement key "key"`,
		},
		{
			name:        "empty input",
			replacement: []string{},
			expected:    map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := getReplacementMap(tt.replacement)
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
			if len(m) != len(tt.expected) {
				t.Fatalf("expected %d entries, got %d: %v", len(tt.expected), len(m), m)
			}
			for k, v := range tt.expected {
				if m[k] != v {
					t.Errorf("expected %q to map to %q, got %q", k, v, m[k])
				}
			}
		})
	}
}

func TestLoadConfigWithReplacements(t *testing.T) {
	tests := []struct {
		name          string
		configOption  configOption
		wantErr       bool
		errContains   string
		expectSheets  int
		expectReplace map[string]string
	}{
		{
			name:         "missing config path",
			configOption: configOption{},
			wantErr:      true,
			errContains:  "configuration file is not specified",
		},
		{
			name: "non-existent file",
			configOption: configOption{
				configFilePath: "/nonexistent/path/config.yml",
			},
			wantErr: true,
		},
		{
			name: "invalid data_type",
			configOption: configOption{
				configFilePath: "testdata/invalid_column_type.yml",
			},
			wantErr:     true,
			errContains: "unsupported data_type [number]",
		},
		{
			name: "valid config with replacements",
			configOption: configOption{
				configFilePath: "testdata/valid_config.yml",
				replacements:   []string{"name:test"},
			},
			expectSheets:  1,
			expectReplace: map[string]string{"name": "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, replacements, err := loadConfigWithReplacements(tt.configOption)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(config.Sheets) != tt.expectSheets {
				t.Errorf("expected %d sheets, got %d", tt.expectSheets, len(config.Sheets))
			}
			if len(replacements) != len(tt.expectReplace) {
				t.Fatalf("expected %d replacements, got %d", len(tt.expectReplace), len(replacements))
			}
			for k, v := range tt.expectReplace {
				if replacements[k] != v {
					t.Errorf("expected %q to map to %q, got %q", k, v, replacements[k])
				}
			}
		})
	}
}
