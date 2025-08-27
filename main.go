package main

import (
	"github.com/spf13/viper"
	"gpt4omini/cmd"
)

func main() {
	viper.SetEnvPrefix(cmd.CliName)
	viper.AutomaticEnv()
	cmd.Execute()
}
