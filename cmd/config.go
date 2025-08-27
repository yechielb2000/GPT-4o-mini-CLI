package cmd

import "github.com/spf13/cobra"

const ConfigCmd = "config"

var configCmd = &cobra.Command{
	Use:   ConfigCmd,
	Short: "View or update config file",
	Long: `Manage the configuration stored in config.yaml.
You can print the current config or update API/model fields.`,
}

func init() {}
