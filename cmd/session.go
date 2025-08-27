package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/session"
	"os"
	"strings"
)

const SessionName string = "session"

var (
	showList        bool
	startNewSession bool
	sessionId       string
	sessionsManager = session.GetSessionsManager()
)

var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  "Manage realtime sessions (create, list, resume, delete) in an interactive CLI until you exit.",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("\n(session-cli) > ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "" {
				continue
			}

			args := strings.Split(input, " ")
			command := args[0]

			handleCommand(command, args[1:])
		}
	},
}

func handleCommand(command string, args []string) {
	switch command {
	case "list":
		for id, s := range sessionsManager.Sessions() {
			fmt.Printf("ID: %s | Type: %s\n", id, s.GetType())
		}
	case "exit", "quit":
		fmt.Println("Exiting session manager...")
		return
	}
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
