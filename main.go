package main

import (
	"github.com/spf13/viper"
	"gpt4omini/cmd"
	"gpt4omini/session"
)

func main() {
	viper.SetEnvPrefix(cmd.CliName)
	viper.AutomaticEnv()
	_ = session.GetSessionsManager()
	cmd.Execute()
}
