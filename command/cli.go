package command

import (
	"fmt"
	"strings"

	"github.com/alexhokl/sql-export/model"
	"github.com/spf13/cobra"
)

// ManagerCli struct
type ManagerCli struct {
}

// NewManagerCli creates a new manager cli instance
func NewManagerCli() *ManagerCli {
	cli := ManagerCli{}
	return &cli
}

// ShowHelp shows the command help
func (cli *ManagerCli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.HelpFunc()(cmd, args)
	return nil
}

type configOption struct {
	configFilePath string
	replacements   []string
}

// loadConfigWithReplacements parses the configuration file and
// replacement flags shared by all export commands
func loadConfigWithReplacements(opts configOption) (*model.ExportConfig, map[string]string, error) {
	if opts.configFilePath == "" {
		return nil, nil, fmt.Errorf("configuration file is not specified")
	}
	config, errConfig := model.ParseConfig(opts.configFilePath)
	if errConfig != nil {
		return nil, nil, errConfig
	}
	replacements, err := getReplacementMap(opts.replacements)
	if err != nil {
		return nil, nil, err
	}
	return config, replacements, nil
}

// NewManagerCommand returns the main command of this exporter
func NewManagerCommand(cli *ManagerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sql-export",
		Short:        "SQL data exporter",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.ShowHelp(cmd, args)
		},
	}

	AddCommands(cmd, cli)
	return cmd
}

func getReplacementMap(replacements []string) (map[string]string, error) {
	m := make(map[string]string, len(replacements))
	for _, r := range replacements {
		if r == "" {
			continue
		}
		splits := strings.SplitN(r, ":", 2)
		if len(splits) != 2 {
			return nil, fmt.Errorf("invalid replacement %q: expected format key:value", r)
		}
		if _, exists := m[splits[0]]; exists {
			return nil, fmt.Errorf("duplicated replacement key %q", splits[0])
		}
		m[splits[0]] = splits[1]
	}
	return m, nil
}
