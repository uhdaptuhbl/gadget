package main

import (
	"fmt"

	"gadget/halt"
	"gadget/harness"
	"gadget/settings"
	"gadget/teapot/cookiejar"
)

func main() {
	var jar = cookiejar.Builder().New()
	fmt.Printf("%+v\n", jar)
	fmt.Printf("%+v\n", halt.InterruptOptions{})
	fmt.Printf("%+v\n", harness.InvalidValueError{})
	fmt.Printf("%+v\n", settings.FlagFactory{})
	// fmt.Printf("%+v\n", teapot.RequestArgs{})
}
