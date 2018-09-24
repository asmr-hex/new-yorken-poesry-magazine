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
}

func (a *API) GetPoem(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEM")
}

func (a *API) GeneratePoem(rw web.ResponseWriter, req *web.Request) {
	// extracting the poetId path param
	poetId := req.PathParams[API_ID_PATH_PARAM]

	poet := &types.Poet{Id: poetId}

	// we need to read the poet before we can use it.
	err := poet.Read(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// set the execution context...
	poet.ExecContext = &a.Config.ExecContext

	// generate poem
	// TODO (cw|9.23.2018) rate limit this...
	poem, err := poet.GeneratePoem()
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	poemJSON, err := json.Marshal(poem)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(poemJSON)
}
