package main

import (
	"fmt"

	"gadget/halt"
	"gadget/harness"
	"gadget/settings"
	// "gadget/teapot"
)

func main() {
	fmt.Printf("%+v\n", halt.InterruptOptions{})
	fmt.Printf("%+v\n", harness.InvalidValueError{})
	fmt.Printf("%+v\n", settings.FlagFactory{})
	// fmt.Printf("%+v\n", teapot.RequestArgs{})
}
