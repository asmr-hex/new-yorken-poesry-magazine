package core

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
	uuid "github.com/satori/go.uuid"
)

func (a *API) GetPoets(rw web.ResponseWriter, req *web.Request) {
	var (
		poets []*types.Poet
		err   error
	)

	poets, err = types.ReadPoets(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	poetsJSON, err := json.Marshal(poets)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(poetsJSON)

	a.Info("successfully read all poets")
}

// creates a poet for a valid, authenticated user.
//
// this API handler is used by both the webpage dashboard and public
// APIs. a notable difference between this POST handler and others in
// this API is that it only accepts Content-Type: multipart/form-data
// (i.e. we can't give it json and expect it to work). the reason for
// this is that we need to be uploading files here and form-data is the
// easier and most efficient way of doing that(?).
//
// an example curl to this endpoint would look something like this,
//
// curl -X POST \
//      -b "session_token=<your-session-token>" \
//      -F "name=<your-poet-name>" \
//      -F "description=<poet-description>" \
//      -F "language=<poet-language>" \
//      -F "src[]=@path/to/program-file;filename=program" \
//      -F "src[]=@path/to/optional/param-file;filename=parameters" \
//      https://poem.cool/dashboard/poet
//      ( or https://poem.cool/api/v1/poet)
//
// note that in the example above there are a couple conventions that
// the user must follow in order for their file to be properly uploaded.
// first, all files must belong to the src[] form key. this way all files
// are bundled in one nice array. finally, in order for the server to know
// which file is which, the filename of each file must be set to either
// 'program' or 'parameters'. this way, the users can name their files whatever
// they want locally.
//
// TODO (cw|4.27.2018) this function is way too large. consider
// refactoring into smaller pieces.
//
func (a *API) CreatePoet(rw web.ResponseWriter, req *web.Request) {
	var (
		err error
		fds = struct { // file descriptors
			program    *multipart.FileHeader
			parameters *multipart.FileHeader
		}{}
	)

	// anyone who sends this request *must* have a session token in their
	// request header (or cookies?) since only logged in users can create poets.

	// get session token from cookie (maybe use golang CookieJar)
	tokenCookie, err := req.Cookie(SESSION_TOKEN_COOKIE_NAME)
	if err != nil {
		// handle this error
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	token := tokenCookie.Value

	// get username from token
	userId, validToken := a.Sessions.GetUserByToken(token)
	if !validToken {
		err = fmt.Errorf("invalid session token!")

		// handle this error
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusUnauthorized)

		return
	}

	// parse multipart-form from request
	err = req.ParseMultipartForm(30 << 20)
	if err != nil {
		switch {
		case err == http.ErrNotMultipart:
			fallthrough
		case err == multipart.ErrMessageTooLarge:
			// handle this error
			a.Error("User Error: %s", err.Error())

			http.Error(rw, err.Error(), http.StatusBadRequest)
		default:
			// log internal error
			a.Error("Internal Error: %s", err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	// iterate over form files
	formFiles := req.MultipartForm.File
	for filesKey, files := range formFiles {
		// we onlye care about the POET_FILES_FORM_KEY
		if filesKey != POET_FILES_FORM_KEY {
			err = fmt.Errorf(
				"Encountered abnormal key, %s",
				filesKey,
			)
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusBadRequest)

			return
		}

		// there should be at most two files
		nFiles := len(files)
		if nFiles > 2 || nFiles < 1 {
			err = fmt.Errorf(
				"Expected at most 2 files within %s form array, given %d",
				POET_FILES_FORM_KEY,
				nFiles,
			)

			a.Error("User Error: %s", err.Error())

			http.Error(rw, err.Error(), http.StatusBadRequest)

			return
		}

		// try to get code files and the optional parameters file
		for _, file := range files {
			switch file.Filename {
			case POET_PROG_FILENAME:
				if fds.program != nil {
					// this means multiple program files were uploaded!
					err = fmt.Errorf("Multiple program files uploaded, only 1 allowed!")
					a.Error("User Error: %s", err.Error())

					http.Error(rw, err.Error(), http.StatusBadRequest)

					return
				}

				fds.program = file

			case POET_PARAMS_FILENAME:
				if fds.parameters != nil {
					// this means multiple parameter files were uploaded!
					err = fmt.Errorf("Multiple parameter files uploaded, only 1 allowed!")
					a.Error("User Error: %s", err.Error())

					http.Error(rw, err.Error(), http.StatusBadRequest)

					return
				}

				fds.parameters = file

			default:
				// invalid filename was included
				err = fmt.Errorf("Invalid filename provided, %s", file.Filename)

				a.Error("User Error: %s", err.Error())

				http.Error(rw, err.Error(), http.StatusBadRequest)

				return
			}
		} // end for
	} // end for

	// ensure that we have a program file
	if fds.program == nil {
		err = fmt.Errorf("No program file was uploaded! At least 1 required.")
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// open up the program file!
	fdProg, err := fds.program.Open()
	defer fdProg.Close()
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// create new poet
	poetID := uuid.NewV4().String()

	// initialize poet struct
	poet := &types.Poet{
		Id:              poetID,
		Designer:        &types.User{Id: userId},
		BirthDate:       time.Now(),
		Name:            req.PostFormValue(POET_NAME_PARAM),
		Description:     req.PostFormValue(POET_DESCRIPTION_PARAM),
		Language:        req.PostFormValue(POET_LANGUAGE_PARAM),
		ProgramFileName: POET_PROG_FILENAME,
		Path:            path.Join(POET_DIR, poetID),
	}

	// validate the poet structure
	err = poet.Validate(
		consts.CREATE,
		types.PoetValidationParams{
			// ensure that the current logged in/authed user
			// can only create a poet for themselves.
			DesignerId: userId,
			SupportedLangs: map[string]bool{
				// TODO (cw|4.27.2018) get this list of
				// supported langauges from somewhere.
				"python": true,
			},
		},
	)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// create new poet directory
	err = os.Mkdir(poet.Path, os.ModePerm)
	if err != nil {
		a.Error(err.Error())

		// returrn response
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// create program file on fs
	dstProg, err := os.Create(filepath.Join(poet.Path, fds.program.Filename))
	defer dstProg.Close()
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// persist program file to the fs
	if _, err = io.Copy(dstProg, fdProg); err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// persist parameters file on disk if provided
	if fds.parameters != nil {
		// open up the parameteres file!
		fdParam, err := fds.parameters.Open()
		defer fdParam.Close()
		if err != nil {
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		// create parameters file on the fs
		dstParam, err := os.Create(filepath.Join(poet.Path, fds.parameters.Filename))
		defer dstParam.Close()
		if err != nil {
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		// persist params file to the fs
		if _, err = io.Copy(dstParam, fdParam); err != nil {
			a.Error(err.Error())

			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		// add this parameters file to the poet
		poet.ParameterFileName = POET_PARAMS_FILENAME
		poet.ParameterFileIncluded = true
	}

	// create poet in db
	err = poet.Create(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// write to the response!
	poet.Sanitize()

	// write json encoded data into response
	err = json.NewEncoder(rw).Encode(poet)
	if err != nil {
		a.Error(err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	a.Info("Poet successfully created ^-^")

	// set execution context for poet
	poet.ExecContext = &a.Config.ExecContext

	a.Info("Testing Poet, %s", poet.Name)
	err = poet.TestPoet()
	if err != nil {
		a.Error(err.Error())
	}
}

func (*API) GetPoet(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POET")
}

func (*API) UpdatePoet(rw web.ResponseWriter, req *web.Request) {
	// we can only update certain fields, like description and avatar
	fmt.Println("TODO UPDATE POET")
}

func (*API) DeletePoet(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO DELETE POET")
}
