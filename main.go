package main

import (
	"os"

	"github.com/vizv/ipfilter/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
