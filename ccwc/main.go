package main

import (
	"fmt"
	"os"

	"github.com/suhalvemu/wc/cmd"
)

var version = "dev" //overridden by Makefile

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("ccwc version", version)
		return
	}
	cmd.Execute()
}
