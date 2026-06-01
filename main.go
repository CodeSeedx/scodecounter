package main

import (
	"github.com/CodeSeedx/scodecounter/cmd"
)

// version is set at build time via -ldflags "-X main.version=xxx"
var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
