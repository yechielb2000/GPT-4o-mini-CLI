package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/session"
	"log"
	"os"
	"strings"
)

const SessionName string = "session"

var sessionsManager = session.GetSessionsManager()

var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  "Manage realtime sessions (create, list, resume, delete) in an interactive CLI until you exit.",
	Run: func(cmd *cobra.Command, args []string) {
		var reader = bufio.NewReader(os.Stdin)
		var exit bool
		for !exit {
			fmt.Print("\n(session-cli) > ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "" {
				return
			}

			args = strings.Split(input, " ")
			command := args[0]
			args = args[1:]

			switch command {
			case "list":
				for id, s := range sessionsManager.Sessions() {
					fmt.Printf("ID: %s | Type: %s\n", id, s.GetType())
				}
			case "new":
				fmt.Println(fmt.Sprintf("Enter session type (%s):", session.GetSessionTypes()))
				chosenType, _ := reader.ReadString('\n')
				chosenType = strings.TrimSpace(chosenType)
				newSession, err := session.NewSessionByType(chosenType)
				if err != nil {
					log.Println("Got an error while trying to create session:", err)
					return
				}
				sessionsManager.AddSession(newSession)
				fmt.Printf("Started new %s session with ID %s\n\n", newSession.GetType(), newSession.GetID())
				newSession.Start()
			case "resume":
				if len(args) < 1 {
					fmt.Println("Usage: resume <sessionID>")
					return
				}
				id := args[0]
				s, err := sessionsManager.GetSession(id)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Printf("Resuming %s session %s\n", s.GetType(), id)
				if err := s.Resume(); err != nil {
					log.Println(err)
				}
			case "delete":
				if len(args) < 1 {
					fmt.Println("Usage: delete <sessionID>")
					return
				}
				id := args[0]
				if err := sessionsManager.RemoveSession(id); err != nil {
					log.Println(err)
				} else {
					fmt.Println("Deleted session:", id)
				}
			case "exit", "quit":
				fmt.Println("Exiting session manager...")
				exit = true
			default:
				fmt.Println("Unknown command. Available: list, new, resume <id>, delete <id>, exit")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
