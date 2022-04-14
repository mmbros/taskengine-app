package main

import (
	"os"

	"github.com/mmbros/taskengine-app/cmd"
)

func main() {
	code := cmd.Execute(os.Stdout)
	os.Exit(code)
}
