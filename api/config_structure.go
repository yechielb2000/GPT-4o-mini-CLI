package api

type Api struct {
	Key    string `json:"key"`
	Host   string `json:"host"`
	Schema string `json:"schema"`
}

type Model struct {
	Name        string `json:"name"`
	Instruction string `json:"instruction"`
}

type Config struct {
	Api   Api   `json:"api"`
	Model Model `json:"model"`
}
