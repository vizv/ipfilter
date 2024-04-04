package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/vizv/ipfilter/cmd"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
