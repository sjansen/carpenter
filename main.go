package main

import (
	"fmt"
	"math"
	"os"

	"github.com/google/skylark"
	"github.com/google/skylark/resolve"
)

const script = `
x = sqrt(float(3*3 + 4*4))
`

func sqrt(
	thread *skylark.Thread,
	_ *skylark.Builtin,
	args skylark.Tuple,
	kwargs []skylark.Tuple,
) (skylark.Value, error) {
	var x skylark.Float
	if err := skylark.UnpackArgs("sqrt", args, kwargs, "x", &x); err != nil {
		return nil, err
	}
	result := math.Sqrt(float64(x))
	return skylark.Float(result), nil
}

func main() {
	resolve.AllowFloat = true

	globals := skylark.StringDict{
		"sqrt": skylark.NewBuiltin("sqrt", sqrt),
	}
	thread := &skylark.Thread{}
	if result, err := skylark.ExecFile(thread, "<stdin>", script, globals); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(result["x"])
	}
}
