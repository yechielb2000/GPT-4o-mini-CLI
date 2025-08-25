package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/session"
	"net/url"
)

const SessionName string = "session"

var (
	showList        bool
	startNewSession bool
	getSession      string
)

/*
sessionCmd subcommand handle all session actions
session -n "start new session"
session --list "get last active sessions"
session "session name" return to session
*/
var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long: fmt.Sprintf(
		`Make session actions using this subcommand.
For example: %s %s -n. This starts a new session`, CliName, SessionName,
	),
	Run: func(cmd *cobra.Command, args []string) {
		if showList {
			// TODO: return active sessions from sessions manager
		} else if startNewSession {
			//TODO: change this to params that can be changed via cli (use config.yaml maybe)
			u := url.URL{
				Scheme:   "wss",
				Host:     "api.openai.com",
				Path:     "/v1/realtime/sessions",
				RawQuery: "model=gpt-4o-realtime-preview-2024-12-17",
			}
			newSession, err := session.NewSession(u, apiKey)
			if err != nil {
				fmt.Println(err)
				return
			}
			newSession.Start()
		}
	},
}

func init() {
	//TODO: read from env or config file
	apiKey = "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"

	rootCmd.AddCommand(sessionCmd)
	sessionCmd.Flags().BoolVarP(&startNewSession, "new", "n", false, "Create new session.")
	sessionCmd.Flags().BoolVarP(&showList, "list", "l", false, "Get last active sessions.")
	sessionCmd.Flags().StringVarP(&getSession, "get", "g", "", "Return to session if exists.")
	sessionCmd.MarkFlagsMutuallyExclusive("new", "get", "list")
}
