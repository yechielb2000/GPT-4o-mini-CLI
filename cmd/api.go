package cmd

import (
	"github.com/spf13/cobra"
)

const ApiCmd string = "api"
const ApiKeyCmd string = "key"

/*
apiCmd subcommand handles all api related actions
*/
var apiCmd = &cobra.Command{
	Use:   ApiCmd,
	Short: "Make api actions",
	Long: `Make api actions.
For example: set the api key of the gpt-4o-mini`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

/*
apiKeyCmd subcommand handles api key actions such as edit api key or print current key.
api key -"newkey" | set new api key.
api key -p | prints current api key.
*/
var apiKeyCmd = &cobra.Command{
	Use:   ApiKeyCmd,
	Short: "key related actions",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	apiCmd.AddCommand(apiKeyCmd)
	apiKeyCmd.Flags().BoolP("print", "p", false, "Print current api key.")
	apiKeyCmd.Flags().StringVar(&apiKey, "new", "", "Set new api key.")
	rootCmd.AddCommand(apiCmd)
}
