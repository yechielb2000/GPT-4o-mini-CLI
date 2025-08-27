package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/config"
)

const ConfigCmd = "config"

var (
	apiKey      string
	apiHost     string
	apiSchema   string
	modelName   string
	modelInstr  string
	printConfig bool
)

var configCmd = &cobra.Command{
	Use:   ConfigCmd,
	Short: "View or update config file",
	Long: `Manage the configuration stored in config.yaml.
You can print the current config or update API/model fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		if apiKey != "" {
			cfg.Api.Key = apiKey
		}
		if apiHost != "" {
			cfg.Api.Host = apiHost
		}
		if apiSchema != "" {
			cfg.Api.Schema = apiSchema
		}
		if modelName != "" {
			cfg.Model.Name = modelName
		}
		if modelInstr != "" {
			cfg.Model.Instruction = modelInstr
		}
		if printConfig {
			fmt.Println("hello there")
			fmt.Println(cfg)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&printConfig, "print", "p", false, "Print current config.")
	configCmd.Flags().StringVar(&apiKey, "key", "", "Set API key.")
	configCmd.Flags().StringVar(&apiHost, "host", "", "Set API host.")
	configCmd.Flags().StringVar(&apiSchema, "schema", "", "Set API schema (wss/ws/http).")
	configCmd.Flags().StringVar(&modelName, "model", "", "Set model name.")
	configCmd.Flags().StringVar(&modelInstr, "instruction", "", "Set model instruction.")
}
