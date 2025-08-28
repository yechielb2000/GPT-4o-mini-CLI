package global_tools

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
