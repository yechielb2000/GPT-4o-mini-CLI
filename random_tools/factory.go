package random_tools

import "fmt"

var Factory = map[string]func(args map[string]interface{}) any{
	"multiply": func(args map[string]interface{}) any {
		return InputOfTwo(multiply, args)
	},
	"add": func(args map[string]interface{}) any {
		return InputOfTwo(add, args)
	},
	"minus": func(args map[string]interface{}) any {
		return InputOfTwo(minus, args)
	},
}

func CallFunction(name string, args map[string]interface{}) (string, error) {
	if f, ok := Factory[name]; ok {
		result := fmt.Sprintf("%v", f(args))
		return result, nil
	}
	return "", fmt.Errorf("no such function: %s", name)
}

// InputOfTwo used to reduce code for the factory
func InputOfTwo(fn func(float64, float64) float64, args map[string]interface{}) float64 {
	x := args["x"].(float64)
	y := args["y"].(float64)
	return fn(x, y)
}
