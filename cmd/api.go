package cmd

import (
	"github.com/spf13/cobra"
)

const ApiCmd string = "api"
const ApiKeyCmd string = "key"

var (
	printKey bool
)

/*
apiCmd subcommand handles all api related actions
*/
var apiCmd = &cobra.Command{
	Use:   ApiCmd,
	Short: "Make api actions",
	Long: `Make api actions.
For example: set the api key of the gpt-4o-mini`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiBaseURL == "" {
			//TODO: validate it better of course..
			//TODO: log error etc..
			return
		}
	},
}

/*
apiKeyCmd subcommand handles api key actions such as edit api key or print current key.
api key --new "newkey" | set new api key.
api key -p | prints current api key.
*/
var apiKeyCmd = &cobra.Command{
	Use:   ApiKeyCmd,
	Short: "key related actions",
	Run: func(cmd *cobra.Command, args []string) {
		if printKey {
			// TODO: log out key
		} else {
			if apiKey == "" {
				// TODO: log out an error (apikey can not be empty)
				// TODO: see if there is a better way to validate user input
				return
			}
		}
	},
}

func init() {
	apiCmd.AddCommand(apiKeyCmd)
	apiKeyCmd.Flags().BoolVarP(&printKey, "print", "p", false, "Print current api key.")
	apiKeyCmd.Flags().StringVar(&apiKey, "new", "", "Set new api key.")
	apiKeyCmd.MarkFlagsMutuallyExclusive("print", "new")
	rootCmd.AddCommand(apiCmd)
	//TODO: see how to use bind
	apiCmd.Flags().StringVarP(&apiBaseURL, "new", "n", "", "Set new base URL.")
}
