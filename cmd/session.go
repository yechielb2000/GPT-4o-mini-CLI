package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const SessionName string = "session"

var sessionCmd = &cobra.Command{
	Use:   SessionName,
	Short: "Make session actions",
	Long:  fmt.Sprintf(`Make session actions\nFor example: %s %s -n. This starts a new session`, CliName, SessionName),
}
