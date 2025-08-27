# GPT-4o-mini-CLI

Real-time GPT-4o-mini CLI with Function Calling

## Installation Time :)

compile and install the cli using these commands:

```shell
go build
go install
chmod +x gptoncli # you probably won't need it but in case you can't execute, run this.
```

Set a config file called `config.yaml`.

```shell
mkdir -p /etc/gptoncli && touch /etc/gptoncli/config.yaml
```

Then edit the file like this (change the values to your values)

```yaml
api:
  key: your-api-key
  host: the host (api.openai.com)
  schema: the schema (wss)
model:
  name: the model name (gpt-4.1)
  instruction: the initial instructions (you are a rock band assistant...)
```

> _Note_: You can also change the values later using the cli

## Using the `gptoncli`

To see usage of the cli You can run `gptoncli --help`.  
Let's see how it looks like very quickly.  
```yaml
gptoncli:
    session: subcommand for handling session actions.
        - new: create new session.
        - list: list all sessions (sessions dies when we stop using the cli).
        - "session id": will resume the session with the requested session. 
    config: subcommand for handling the config file
        - print: print the current config file
        - key: update the api key.
        - host: update the target host.
        - schema: update the schema type.
        - model: update the model type.
        - instruction: update the default initiative instructions for each session.
```
