package go2js

import (
	"fmt"
	"syscall/js"
)

// invokeJS calls javascript function with the given arguments.
func invokeJS(name string, args ...any) (js.Value, error) {
	if fn, ok := getJS(name); ok {
		return fn.Invoke(args...), nil
	}
	return js.Undefined(), fmt.Errorf("function \"%s\" is undefined", name)
}

// getJS returns a JavaScript function object for a given
// function name, and whether its a defined function or not.
func getJS(name string) (js.Value, bool) {
	val := js.Global().Get(name)
	if !val.IsUndefined() && val.Type() == js.TypeFunction {
		return val, true
	}
	return js.Undefined(), false
}
