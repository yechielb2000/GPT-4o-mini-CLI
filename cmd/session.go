package cmd

import (
	"github.com/spf13/cobra"
	"gpt4omini/session"
)

const SessionName string = "session"

var (
	showList        bool
	startNewSession bool
	sessionId       string
	sessionsManager = session.GetSessionsManager()
)

/*
sessionCmd subcommand handle all session actions
session -new "start new session"
session --list "get last active sessions"
session "session name" return to session
*/
var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
