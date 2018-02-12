package main

import (
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/core"
)

func main() {
	// init platform instance
	p := core.NewPlatform()

	p.Success("New Yorken Poesry daemon running...\n")

	p.Start()
}
