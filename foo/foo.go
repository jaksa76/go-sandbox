package main

import "fmt"

// SayHello says hello
func SayHello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	fmt.Println(SayHello("World"))
}
