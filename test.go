package main

import (
	"fmt"

	"bitbucket.services.cymru.com/fst/gadget/halt"
	"bitbucket.services.cymru.com/fst/gadget/harness"
	"bitbucket.services.cymru.com/fst/gadget/settings"
	"bitbucket.services.cymru.com/fst/gadget/teapot"
)

func main() {
	fmt.Printf("%+v\n", halt.InterruptOptions{})
	fmt.Printf("%+v\n", harness.InvalidValueError{})
	fmt.Printf("%+v\n", settings.FlagFactory{})
	fmt.Printf("%+v\n", teapot.RequestArgs{})
}
