package global_tools

func InputOfTwo(fn func(float64, float64) float64, args map[string]interface{}) float64 {
	x := args["x"].(float64)
	y := args["y"].(float64)
	return fn(x, y)
}

func multiply(x float64, y float64) float64 {
	return x * y
}

func add(x float64, y float64) float64 {
	return x + y
}

func minus(x float64, y float64) float64 {
	return x - y
}
