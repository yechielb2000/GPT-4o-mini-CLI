package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const SessionName string = "session"

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
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			// TODO: return active sessions from sessions manager
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.Flags().BoolP("new", "n", false, "Create new session.")
	sessionCmd.Flags().BoolP("list", "l", false, "Get last active sessions.")
	sessionCmd.Flags().StringP("get", "g", "", "Return to session if exists.")
	sessionCmd.MarkFlagsMutuallyExclusive("new", "get", "list")
}
