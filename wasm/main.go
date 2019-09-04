package main

import (
	"syscall/js"
)

var (
	buffer [800 * 600 * 4]byte
	a      byte
)

func getBufferPtr(this js.Value, inputs []js.Value) interface{} {
	return &buffer
}

func printMessage(this js.Value, inputs []js.Value) interface{} {
	message := inputs[0].String()

	document := js.Global().Get("document")
	p := document.Call("createElement", "p")
	p.Set("innerHTML", message)
	document.Get("body").Call("appendChild", p)

	return nil
}

func main() {
	c := make(chan bool)
	js.Global().Set("printMessage", js.FuncOf(printMessage))
	js.Global().Set("buff", js.FuncOf(getBufferPtr))
	<-c
}
