package core

import (
	"fmt"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	"github.com/gocraft/web"
)

// Validate incoming requests.
// ensure path parameters are valid, etc.
func (a *API) Validate(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	fmt.Println(req.RequestURI)

	// validate path parameter strings
	err := ValidateParams(req.PathParams)
	if err != nil {
		// return an error to client
	}

	next(rw, req)
}

func ValidateParams(params map[string]string) error {
	for k, v := range params {
		// apply param specific validation, e.g. id params must be uuid v4, etc.
		switch k {
		case API_ID_PATH_PARAM:
			// id path params MUST be V4 UUIDs
			if !utils.IsValidUUIDV4(v) {
				return fmt.Errorf("Id parameter must be a UUID V4 (given %s)", v)
			}
		default:
			// hmmm, we should never get here.
			return fmt.Errorf("Unknown path parameter %s", k)
		}
	}

	return nil
}

// Authorize requests.
// ensure a user cannot delete other users, etc.
func (*API) Authorize(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// TODO

	next(rw, req)
}
