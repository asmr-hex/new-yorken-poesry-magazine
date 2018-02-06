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

	// init API routers + state
	api := core.NewAPI()

	fmt.Printf(color.GreenString("New Yorken Poesry daemon running...\n"))

	err := http.ListenAndServe("localhost:3000", api.Router)
	if err != nil {
		log.Fatal(err)
	}
}
