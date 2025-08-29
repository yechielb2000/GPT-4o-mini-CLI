# Real-time GPT-4o-mini CLI

A terminal-based client for interacting with OpenAI’s Realtime API.  
This tool supports user input, function calls, and displays model responses in an interactive, live-updating terminal
interface.

## Installation & Configuration

```shell
git clone https://github.com/yechielb2000/GPT-4o-mini-CLI.git
cd GPT-4o-mini-CLI
go build; go install
chmod +x gpt4omini
./gpt4omini
```

> _Note_: You should have a config.yaml under the same directory with your binary.

You can configure the model and the api you are using via this `config.yaml`.  
The structure of the config yaml is [here](config/config_structure.go).

To configure via cli, you can run:

```shell
./gpt4omini config --help
```

You can set the `model` and the `instructions` for specific session but the default is taken from the
config file.

## Basic Usage

To start a session, You can run:

```shell
./gpt4omini session -t 'session-type'
```

This will create a conversation with the model you chose and with the session type you wanted.

To list all sessions types, you can run:

```shell
./gpt4omini session -s
```

For all flags and usage you can run:

````shell
./gpt4omini session --help
````

