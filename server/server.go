package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/core"
	"github.com/fatih/color"
)

var (
	arg = flag.String("arg", "123", "ok")
)

func main() {
	flag.Parse()

	// check args

	// init platform instance
	p := core.NewPlatform()

	fmt.Printf(color.GreenString("New Yorken Poesry daemon running...\n"))

	err := http.ListenAndServe("0.0.0.0:8080", p.Api.Router)
	if err != nil {
		log.Fatal(err)
	}
}
