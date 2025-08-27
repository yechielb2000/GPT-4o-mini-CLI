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

var (
	showList        bool
	startNewSession bool
	sessionId       string
	sessionsManager = session.GetSessionsManager()
	reader          = bufio.NewReader(os.Stdin)
)

var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  "Manage realtime sessions (create, list, resume, delete) in an interactive CLI until you exit.",
	Run: func(cmd *cobra.Command, args []string) {

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
		go newSession.Start()
	case "resume":
		if len(args) < 2 {
			fmt.Println("Usage: resume <sessionID>")
			return
		}
		id := args[1]
		s, err := sessionsManager.GetSession(id)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Resuming %s session %s\n", s.GetType(), id)
		go s.Start()
	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: delete <sessionID>")
			return
		}
		id := args[1]
		sessionsManager.RemoveSession(id)
		fmt.Println("Deleted session:", id)
	case "exit", "quit":
		fmt.Println("Exiting session manager...")
		return
	}
}

func init() {
	rootCmd.AddCommand(sessionCmd)
}
