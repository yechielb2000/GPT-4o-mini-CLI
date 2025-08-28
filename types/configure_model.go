package types

type ConfigureModelRequest struct {
	Modalities   []string `json:"modalities"`
	Model        string   `json:"model"`
	Instructions string   `json:"instructions"`
	Tools        []Tool   `json:"tools"`
}

type ConfigureModelResponse struct {
	Id           string       `json:"id"`
	Object       string       `json:"object"`
	ClientSecret ClientSecret `json:"client_secret"`
}
