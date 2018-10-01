package core

import (
	"encoding/json"
	"net/http"

	"github.com/frenata/xaqt"
	"github.com/gocraft/web"
)

func (a *API) GetSupportedLanguages(rw web.ResponseWriter, req *web.Request) {
	languages, err := a.getSupportedLanguages()
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	languagesJSON, err := json.Marshal(languages)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(languagesJSON)
}

func (*API) getSupportedLanguages() ([]string, error) {
	ctx, err := xaqt.NewContext(xaqt.GetCompilers())
	if err != nil {
		return nil, err
	}

	// TODO (cw|9.25.2018) return more sophisticated data about languages including
	// version and supported libraries...
	languages := []string{}
	for k, _ := range ctx.Languages() {
		languages = append(languages, k)
	}

	return languages, nil
}
