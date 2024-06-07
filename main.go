package main

import (
	"fmt"

	"github.com/portilho13/vcs-cli-go/args"
)

func main() {
	fmt.Println("Hello, World!")
	commandArgs := args.GetArgs()
	switch commandArgs[0] {
		case "init":
	}
}