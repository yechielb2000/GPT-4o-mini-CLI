package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/session"
	"log"
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
	Long: fmt.Sprintf(
		`Make session actions using this subcommand.
For example: %s %s -n. This starts a new session`, CliName, SessionName,
	),
	Run: func(cmd *cobra.Command, args []string) {
		if showList {
			for id := range sessionsManager.Sessions() {
				fmt.Println("RealtimeSession ID:", id)
			}
		} else if startNewSession {
			newSession, err := session.NewRealtimeSession()
			if err != nil {
				log.Fatal(err)
			}
			sessionsManager.AddSession(newSession)
			if err != nil {
				fmt.Println(err)
				return
			}
			newSession.Start()
		} else if sessionId != "" {
			s, err := sessionsManager.GetSession(sessionId)
			if err != nil {
				log.Print(err)
			}
			s.Start()
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.Flags().BoolVarP(&startNewSession, "new", "n", false, "Create new session.")
	sessionCmd.Flags().BoolVarP(&showList, "list", "l", false, "Get last active sessions.")
	sessionCmd.Flags().StringVarP(&sessionId, "get", "g", "", "Return to session if exists.")
	sessionCmd.MarkFlagsMutuallyExclusive("new", "get", "list")
}
