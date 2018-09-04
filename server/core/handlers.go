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

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	"github.com/gocraft/web"
	uuid "github.com/satori/go.uuid"
)

/*

  Plural Types

*/
func (a *API) GetUsers(rw web.ResponseWriter, req *web.Request) {
	var (
		users []*types.User
		err   error
	)

	users, err = types.ReadUsers(a.db)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		a.Error(err.Error())

		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(usersJSON)

	a.Info("successfully read all users")
}

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
func (a *API) SignUp(rw web.ResponseWriter, req *web.Request) {
	// this should be a request handler for a registration endpoint

	// get data from request

	// create registration token

	// store registration token and data in db

	// send email with a link which will have the token in the url
	// and the page you are directed to will send a POST request with
	// the provided token which will then fetch the corresponding user
	// data temporarily keyed by the token and hit the CreateUser hanlder.

	// TODO (cw|4.25.2018) everything below here will go into the VerifyAccount
	// handler.

	// For now this will *only* create a new user.
	a.CreateUser(rw, req)
}

// verifies that the account about to be created has been associated with and
// email address.
//
// TODO (cw|4.25.2018) get a verification key from the request param, lookup the
// user about to be created within an in memory cache, if it exists then proceed,
// if not, then this is an unverified email address.
//
func (a *API) VerifyAccount(rw web.ResponseWriter, req *web.Request) {
	// TODO once we implement email verification tokens, we will be expecting
	// a registration token in the data, which we will here use to fetch the
	// appropriate data from the db. Once we have the user, we will proceed as
	// normally below. ✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。✲´*。.❄¨¯`*✲。❄。*。
}

// creates a user account.
//
// this will invariably be done after signup and email verification has occured.
//
// upon creation, a user is immediately logged in. consequently, CreateUser should
// never write to the response body so Login is able to.
//
// TODO (cw|4.25.2018) currently there is no distinction between signup and create
// user, but eventually signup should required email verification. When this occurs,
// this handler should no longer be a top-level handler, but called from the
// VerifyAccount func which will pass in the relevant info about the user about to
// be created (that info will not be in the request but cached in memory on SignUp).
// For that reason, we will eventually need to change the signature of this function
// to accept a User struct and a ResponseWriter (we no longer need the Request).
//
func (a *API) CreateUser(rw web.ResponseWriter, req *web.Request) {
	var (
		user types.User
		err  error
	)

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&user)
	if err != nil {
		a.Error("Unable to decode POST raw-data: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	err = user.Validate(consts.CREATE)
	if err != nil {
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// assign user id
	user.Id = uuid.NewV4().String()

	// insert data into db tables
	err = user.Create(user.Id, a.db)
	if err != nil {
		a.Error(err.Error())

		// send response
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	a.Info("user %s successfully created!", user.Username)

	// after user is created we can then immediately log the user in
	a.login(&user, rw)
}

// handles login requests.
//
// this handler services requests from both the public api and the
// webapp dashboard. the inbound request should be have `Method: POST`
// with the credentials data json encoded. Here is an example curl,
//
// curl \
//  -X POST \
//  -H "Content-Type: application/json" \
//  -d '{"username": "percival", "password": "phantoms moving mistily"}'
//  https://poem.cool/dashboard/login
//  ( or https://poem.cool/api/v1/login)
//
// upon successful authentication, a unique session token (with an
// expiration date) is assigned to the user if not done so already. This
// token is provided in the response as a cookie, session_token.
//
// finally, the existing/authenticated user is logged in.
//
func (a *API) Login(rw web.ResponseWriter, req *web.Request) {
	var (
		user types.User
		err  error
	)

	defer req.Body.Close()

	// decode body data into a User struct
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&user)
	if err != nil {
		a.Error("Unable to decode POST raw-data: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	err = user.Validate(consts.LOGIN)
	if err != nil {
		a.Error("User Error: %s", err.Error())

		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// authenticate user with password
	err = user.Authenticate(a.db)
	if err != nil {
		a.Error("User Error: %s", err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	a.login(&user, rw)
}

// logs in an existing/authenticated user.
//
// given a User and a ResponseWriter, fetch a session token for this user,
// write the session token and a cookie in the response, write the user json
// data to the response body.
//
// TODO (cw|4.25.2018) if the request is not from a browser
// (i.e. User-Agent is blank or something), the we should include the
// session token within the response payload?
//
func (a *API) login(user *types.User, rw web.ResponseWriter) {
	var (
		err error
	)

	// get session token
	sessionToken := a.Sessions.GetTokenByUser(user.Id)

	// set the sessionToken within a response cookie
	http.SetCookie(rw, &http.Cookie{
		Name:  SESSION_TOKEN_COOKIE_NAME,
		Value: sessionToken,
	})

	// read full user data
	err = user.Read(a.db)
	if err != nil {
		a.Error(err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	user.Sanitize()

	// write json encoded data into response
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		a.Error(err.Error())

		// return response
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	a.Info("user %s successfully logged in!", user.Username)
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
		Designer:        userId,
		Name:            req.PostFormValue(POET_NAME_PARAM),
		Description:     req.PostFormValue(POET_DESCRIPTION_PARAM),
		Language:        req.PostFormValue(POET_LANGUAGE_PARAM),
		ProgramFileName: POET_PROG_FILENAME,
		ExecPath:        path.Join(POET_DIR, poetID),
	}

	// validate the poet structure
	err = poet.Validate(
		consts.CREATE,
		types.PoetValidationParams{
			// ensure that the current logged in/authed user
			// can only create a poet for themselves.
			Designer: userId,
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
	err = os.Mkdir(poet.ExecPath, os.ModePerm)
	if err != nil {
		a.Error(err.Error())

		// returrn response
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// create program file on fs
	dstProg, err := os.Create(filepath.Join(poet.ExecPath, fds.program.Filename))
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
		dstParam, err := os.Create(filepath.Join(poet.ExecPath, fds.parameters.Filename))
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
	err = poet.Create(poetID, a.db)
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
