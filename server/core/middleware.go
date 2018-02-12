package core

import (
	"fmt"
	"regexp"

	"github.com/gocraft/web"
)

var (
	uuidRegexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
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
			if !IsValidUUIDV4(v) {
				return fmt.Errorf("Id parameter must be a UUID V4 (given %s)", v)
			}
		default:
			// hmmm, we should never get here.
			return fmt.Errorf("Unknown path parameter %s", k)
		}
	}

	return nil
}

// check to see if a string is a UUID V4
func IsValidUUIDV4(uuid string) bool {
	return uuidRegexp.MatchString(uuid)
}

// Authorize requests.
// ensure a user cannot delete other users, etc.
func (*API) Authorize(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// TODO

	next(rw, req)
}
