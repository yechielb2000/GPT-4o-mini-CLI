package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

const CliName string = "gpt4omini"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   CliName,
	Short: "Real-time GPT-4o-mini CLI with Function Calling",
	Long: `This is a Command Line Interface (CLI) tool in Go that
interacts with OpenAI’s GPT-4o-realtime-mini, with a websocket, in real time.
The CLI allows users to send messages and receive responses in a streaming format,
simulating the experience of a real-time conversation. Additionally, the CLI supports
function calling, for example: implementing a simple function that multiplies two numbers.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gpt4omini.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
