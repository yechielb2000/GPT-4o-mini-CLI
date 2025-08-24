package main

import (
	"github.com/spf13/viper"
	"gpt4omini/cmd"
)

/*
I want to be able to
- create a new session
- return to one of the last sessions
- get / update api key

simulating commands
session -n "start new session"
session --list "get last active sessions"
session "session name" return to session
api key -u "newapikey"
api key (print new key)
*/

func main() {
	viper.SetEnvPrefix(cmd.CliName)
	viper.AutomaticEnv()
	cmd.Execute()
}
