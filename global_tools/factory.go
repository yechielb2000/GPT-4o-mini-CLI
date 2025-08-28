package global_tools

import (
	"encoding/json"
	"log"
)

type FunctionCallType struct {
	Name string
	Args map[string]interface{}
}

var Factory = map[string]func(args map[string]interface{}) any{
	"multiply": func(args map[string]interface{}) any {
		x := args["x"].(float64)
		y := args["y"].(float64)
		return multiply(x, y)
	},
	"add": func(args map[string]interface{}) any {
		x := args["x"].(float64)
		y := args["y"].(float64)
		return add(x, y)
	},
}

func GetFunctionCallFromItem(item map[string]interface{}) (FunctionCallType, error) {
	functionCall := FunctionCallType{
		Name: item["name"].(string),
	}
	var args map[string]interface{}
	switch v := item["arguments"].(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &args); err != nil {
			log.Println("failed to unmarshal arguments:", err)
			return functionCall, err
		}
	case map[string]interface{}:
		args = v
	}
	functionCall.Args = args
	return functionCall, nil
}
