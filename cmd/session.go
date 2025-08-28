package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpt4omini/config"
	"gpt4omini/session"
	"log"
)

const SessionName string = "session"

var (
	sessionType string
	showTypes   bool
	instruction string
	model       string
)

var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  "Create session on given session type",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		if sessionType != "" {
			if instruction != "" {
				cfg.Model.Instruction = instruction
			}
			if model != "" {
				cfg.Model.Name = model
			}
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
	sessionCmd.Flags().StringVarP(&instruction, "instruction", "i", "", "Instruction to use, default from config file")
	sessionCmd.Flags().StringVarP(&model, "model", "m", "", "Model to use, default from config file")
	sessionCmd.MarkFlagsMutuallyExclusive("type", "show")
}
