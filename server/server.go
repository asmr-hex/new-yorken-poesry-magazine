package main

import (
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/core"
)

// wow this is surprisingly terse (＾▽＾)
func main() {
	// init platform instance
	p := core.NewPlatform()

	p.Success("New Yorken Poesry daemon running...\n")

	p.Start()
}
