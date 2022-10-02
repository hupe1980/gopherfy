package main

import (
	"github.com/hupe1980/gopherfy/cmd"
)

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
