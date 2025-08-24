package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ApiCmd string = "api"

var apiCmd = &cobra.Command{
	Use:   ApiCmd,
	Short: "Make api actions",
	Long: `Make api actions.
For example: set the api key of the gpt-4o-mini`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiKeyName := "apikey"
	viper.SetDefault(apiKeyName, viper.GetString("GPT4oMINI_APIKEY"))
	// should i assign it to rootCmd ?
	apiCmd.PersistentFlags().String(apiKeyName, "", "An API key for the sessions")
	if err := rootCmd.MarkFlagRequired(apiKeyName); err != nil {
		//TODO: log out error and exit
		return
	}
	if err := viper.BindPFlag(apiKeyName, rootCmd.PersistentFlags().Lookup(apiKeyName)); err != nil {
		//TODO: log out error and exit
		return
	}
}
