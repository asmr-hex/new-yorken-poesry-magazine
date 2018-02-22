package core

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
	uuid "github.com/satori/go.uuid"
)

/*

  Plural Types

*/
func (*API) GetUsers(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET USERS")
}

func (*API) GetPoets(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POETS")
}

func (*API) GetPoems(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEMS")
}

func (*API) GetIssues(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET ISSUES")
}

func (*API) GetCommittees(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET COMMITTEES")
}

/*

  User CRUD

*/
// TODO send email verification
func (a *API) RegisterUser(rw web.ResponseWriter, req *web.Request) {
	// this should be a request handler for a registration endpoint

	// get data from request

	// create registration token

	// store registration token and data in db

	// send email with a link which will have the token in the url
	// and the page you are directed to will send a POST request with
	// the provided token which will then fetch the corresponding user
	// data temporarily keyed by the token and hit the CreateUser hanlder.
}

func (a *API) CreateUser(rw web.ResponseWriter, req *web.Request) {
	var (
		user types.User
		err  error
	)

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&user)
	if err != nil {
		a.Error("Unable to decode POST raw-data")

		// TODO send failure response to client
	}

	err = user.Validate(consts.CREATE)
	if err != nil {
		a.Error(err.Error())
	}

	// TODO once we implement email verification tokens, we will be expecting
	// a registration token in the data, which we will here use to fetch the
	// appropriate data from the db. Once we have the user, we will proceed as
	// normally below. ✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。

	// assign user id
	user.Id = uuid.NewV4().String()

	// insert data into db tables
	err = user.Create(user.Id, a.db)
	if err != nil {
		a.Error(err.Error())

		// send response

		return
	}

	// send success response

	a.Info("user %s successfully created!", user.Username)
}

func (a *API) Login(rw web.ResponseWriter, req *web.Request) {
	var (
		err error
	)

	err = req.ParseMultipartForm(30 << 20)
	if err != nil {
		a.Error(err.Error())

		// TODO send response

		return
	}

	fmt.Println(req.Form)

	// get the username, password from request
	user := &types.User{
		Username: req.PostFormValue(LOGIN_USERNAME_PARAM),
		Password: req.PostFormValue(LOGIN_PASSWORD_PARAM),
	}

	fmt.Println(user)

	err = user.Validate(consts.LOGIN)
	if err != nil {
		a.Error(err.Error())

		// send response

		return
	}

	// authenticate user with password
	err = user.Authenticate(a.db)
	if err != nil {
		a.Error(err.Error())

		// return response

		return
	}

	// get session token
	sessionToken := a.Sessions.GetTokenByUser(user.Id)

	a.Info("user %s successfully logged in!", user.Username)
	// return response with session token
	fmt.Println(sessionToken)

	// TODO send successful response WITH SESSION TOKEN IN COOKIES
}

func (a *API) GetUser(rw web.ResponseWriter, req *web.Request) {
	var (
		err error
	)

	// extracting the id path param
	id := req.PathParams[API_ID_PATH_PARAM]

	// assigning said id to id of user struct
	user := &types.User{Id: id}
	err = user.Validate(consts.READ)
	if err != nil {
		a.Error(err.Error())
	}

	// invoke read
	err = user.Read(a.db)
	if err != nil {
		a.Error(err.Error())
	}

	// send actual response back to users

	fmt.Println("TODO GET USER")
}

func (*API) UpdateUser(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO UPDATE USER")
}

func (*API) DeleteUser(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO DELETE USER")
}

/*

  Poet CRD

*/
func (a *API) CreatePoet(rw web.ResponseWriter, req *web.Request) {
	var (
		err error
		fds = struct {
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

		// TODO return response

		return
	}

	token := tokenCookie.Value

	// get username from token
	userId, validToken := a.Sessions.GetUserByToken(token)
	if !validToken {
		err = fmt.Errorf("invalid session token!")

		// handle this error
		a.Error("User Error: %s", err.Error())

		// TODO return response

		return
	}

	// parse multipart-form from request
	err = req.ParseMultipartForm(30 << 20)
	if err != nil {
		// handle this error
		a.Error("User Error: %s", err.Error())

		// TODO return response

		return
	}

	// iterate over form files
	formFiles := req.MultipartForm.File
	for filesKey, files := range formFiles {
		// we onlye care about the POET_FILES_FORM_KEY
		if filesKey != POET_FILES_FORM_KEY {
			a.Error("Encountered abnormal key, %s", filesKey)

			// TODO return http error response

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

			// TODO return response

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
					// TODO return error response

					return
				}

				fds.program = file

			case POET_PARAMS_FILENAME:
				if fds.parameters != nil {
					// this means multiple parameter files were uploaded!
					err = fmt.Errorf("Multiple parameter files uploaded, only 1 allowed!")
					a.Error("User Error: %s", err.Error())
					// TODO return error response

					return
				}

				fds.parameters = file

			default:
				// invalid filename was included
				err = fmt.Errorf("Invalid filename provided, %s", file.Filename)

				a.Error("User Error: %s", err.Error())

				// TODO should we return an error response?

				return
			}
		} // end for
	} // end for

	// ensure that we have a program file
	if fds.program == nil {
		err = fmt.Errorf("No program file was uploaded! At least 1 required.")
		a.Error("User Error: %s", err.Error)

		// TODO return error response

		return
	}

	// open up the program file!
	fdProg, err := fds.program.Open()
	defer fdProg.Close()
	if err != nil {
		a.Error(err.Error())

		// TODO return response

		return
	}

	// create new poet
	poetID := uuid.NewV4().String()

	// initialize poet struct
	poet := &types.Poet{
		Designer:    userId,
		Name:        req.PostFormValue(POET_NAME_PARAM),
		Description: req.PostFormValue(POET_DESCRIPTION_PARAM),
		Language:    req.PostFormValue(POET_LANGUAGE_PARAM),
	}

	// validate the poet structure
	err = poet.Validate(consts.CREATE)
	if err != nil {
		a.Error(err.Error())

		// TODO return responses

		return
	}

	// create new poet directory
	err = os.Mkdir(path.Join(POET_DIR, poetID), os.ModePerm)
	if err != nil {
		a.Error(err.Error())

		// returrn response

		return
	}

	// create program file on fs
	dstProg, err := os.Create(path.Join(POET_DIR, poetID, fds.program.Filename))
	defer dstProg.Close()
	if err != nil {
		a.Error(err.Error())

		// TODO return response (internal server error from http pkg)

		return
	}

	// persist program file to the fs
	if _, err = io.Copy(dstProg, fdProg); err != nil {
		a.Error(err.Error())

		// TODO return response

		return
	}

	// persist parameters file on disk if provided
	if fds.parameters != nil {
		// open up the parameteres file!
		fdParam, err := fds.parameters.Open()
		defer fdParam.Close()
		if err != nil {
			a.Error(err.Error())

			// TODO return response

			return
		}

		// create parameters file on the fs
		dstParam, err := os.Create(path.Join(POET_DIR, poetID, fds.parameters.Filename))
		defer dstParam.Close()
		if err != nil {
			a.Error(err.Error())

			// TODO return response

			return
		}

		// persist params file to the fs
		if _, err = io.Copy(dstParam, fdParam); err != nil {
			a.Error(err.Error())

			// TODO return response

			return
		}
	}

	// create poet in db
	err = poet.Create(poetID, a.db)
	if err != nil {
		a.Error(err.Error())

		// TODO return http response

		return
	}

	a.Info("Poet successfully created ^-^")
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

/*

  Poem Read

*/
func (*API) GetPoem(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET POEM")
}

/*

  Issue Read

*/
func (*API) GetIssue(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET ISSUE")
}

/*

  Committee Read

*/
func (*API) GetCommittee(rw web.ResponseWriter, req *web.Request) {
	fmt.Println("TODO GET COMMITTEE")
}
