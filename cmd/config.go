package cmd

import "github.com/spf13/cobra"

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
}

func init() {
	rootCmd.AddCommand(configCmd)

	apiCmd.Flags().BoolVarP(&printConfig, "print", "p", false, "Print current config.")
	apiCmd.Flags().StringVar(&apiKey, "key", "", "Set API key.")
	apiCmd.Flags().StringVar(&apiHost, "host", "", "Set API host.")
	apiCmd.Flags().StringVar(&apiSchema, "schema", "", "Set API schema (wss/ws/http).")
	apiCmd.Flags().StringVar(&modelName, "model", "", "Set model name.")
	apiCmd.Flags().StringVar(&modelInstr, "instruction", "", "Set model instruction.")
}
