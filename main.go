package main

import "gpt4omini/cmd"

/*
I want to be able to
- create a new session
- return to one of the last sessions
- get / update api key

simulating commands
new session add (session name)
session (session name)
api key -u "newapikey"
api key (print new key)
*/

func main() {
	cmd.Execute()
}
