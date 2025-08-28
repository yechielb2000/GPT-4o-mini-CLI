package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/session"
	"log"
)

const SessionName string = "session"

var (
	sessionType string
	showTypes   bool
)

var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  "Create session on given session type",
	Run: func(cmd *cobra.Command, args []string) {
		if sessionType != "" {
			newSession, err := session.NewSessionByType(sessionType)
			if err != nil {
				log.Println("Got an error while trying to create session:", err)
				return
			}
			fmt.Printf("Started new %s session with ID %s\n\n", newSession.GetType(), newSession.GetID())
			newSession.Start()
		}
		if showTypes {
			fmt.Printf("types: %s", session.GetSessionTypes())
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.Flags().StringVarP(&sessionType, "type", "t", "", "Session type to use")
	sessionCmd.Flags().BoolVarP(&showTypes, "show", "s", false, "List session types available")
	sessionCmd.MarkFlagsMutuallyExclusive("type", "show")
}
