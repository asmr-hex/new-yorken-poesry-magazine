package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/gocraft/web"
)

var (
	arg = flag.String("arg", "123", "ok")
)

type Context struct {
}

func (c *Context) PrintShit(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("OHOHOHHO")
}

func main() {

	flag.Parse()

	// check args

	// define API routes
	router := web.New(Context{}).
		Get("/", (*Context).PrintShit)

	fmt.Printf(color.GreenString("New Yorken Poesry daemon running...\n"))

	err := http.ListenAndServe("localhost:3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
