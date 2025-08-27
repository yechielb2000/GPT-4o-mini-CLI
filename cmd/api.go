package cmd

import (
	"github.com/spf13/cobra"
)

const ApiCmd string = "config"
const ApiKeyCmd string = "key"

var (
	printKey   bool
	apiBaseURL string
)

// apiCmd subcommand handles all config related actions
var apiCmd = &cobra.Command{
	Use:   ApiCmd,
	Short: "Make config actions",
	Long: `Make config actions.
For example: set the config key of the gpt-4o-mini`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiBaseURL == "" {
			//TODO: validate it better of course..
			//TODO: log error etc..
			return
		}
	},
}

/*
apiKeyCmd subcommand handles config key actions such as edit config key or print current key.
config key --new "newkey" | set new config key.
config key -p | prints current config key.
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
	apiKeyCmd.Flags().BoolVarP(&printKey, "print", "p", false, "Print current config key.")
	apiKeyCmd.Flags().StringVar(&apiKey, "new", "", "Set new config key.")
	apiKeyCmd.MarkFlagsMutuallyExclusive("print", "new")
	rootCmd.AddCommand(apiCmd)
	//TODO: see how to use bind
	apiCmd.Flags().StringVarP(&apiBaseURL, "new", "n", "", "Set new base URL.")
}
