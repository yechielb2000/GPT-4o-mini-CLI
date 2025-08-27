package cmd

import (
	"encoding/json"
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
		updated := false
		if apiKey != "" {
			cfg.Api.Key = apiKey
			updated = true
		}
		if apiHost != "" {
			cfg.Api.Host = apiHost
			updated = true
		}
		if apiSchema != "" {
			cfg.Api.Schema = apiSchema
			updated = true
		}
		if modelName != "" {
			cfg.Model.Name = modelName
			updated = true
		}
		if modelInstr != "" {
			cfg.Model.Instruction = modelInstr
			updated = true
		}

		if updated {
			if err := cfg.Save(); err != nil {
				fmt.Println(err)
				return
			}
		}

		if printConfig {
			cfgJson, _ := json.Marshal(cfg)
			fmt.Println(string(cfgJson))
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
