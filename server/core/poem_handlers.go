package core

import (
	"fmt"

	"github.com/gocraft/web"
)

func (*API) GetPoems(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEMS")
}

func (*API) GetPoem(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEM")
}
