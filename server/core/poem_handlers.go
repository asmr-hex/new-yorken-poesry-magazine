package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
)

func (a *API) GetPoems(rw web.ResponseWriter, req *web.Request) {
	var (
		poems []*types.Poem
		err   error
	)

	poems, err = types.ReadPoems(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	poemsJSON, err := json.Marshal(poems)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(poemsJSON)

	a.Info("successfully read all poems")
}

func (a *API) GetPoem(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEM")
}
