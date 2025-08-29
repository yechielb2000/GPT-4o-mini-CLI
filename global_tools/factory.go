package global_tools

import "fmt"

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

func CallFunction(name string, args map[string]interface{}) (string, error) {
	if f, ok := Factory[name]; ok {
		result := fmt.Sprintf("%v", f(args))
		return result, nil
	}
	return "", fmt.Errorf("no such function: %s", name)
}
