package types

import (
	"database/sql"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	_ "github.com/lib/pq"
)

const (
	POET_DESCRIPTION_MAX_CHARS = 2000
)

// Notes about poet executables:
// (1) executables can be named anything, but will be renamed to some standard upon upload
// (2) parameter data can be optionally uploaded if the model decides to store parameters in an external file
// executables will be stored on the filesystem in a safe dir with the path /some/path/bin/<poetId>/

type Poet struct {
	Id          string    `json:"id"`
	Designer    string    `json:"designer"`            // the writer of the poet (user)
	BirthDate   time.Time `json:"birthDate"`           // so we can show years active
	DeathDate   time.Time `json:"deathDate,omitempty"` // this should be set to null for currently active poets
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	ExecPath    string    `json:"-"` // or possibly a Path, this is the path to the source code
	// TODO additional statistics: specifically, it would be cool to see the success rate
	// of a particular poet along with the timeline of how their poems have been recieved

	// what if we also had a poet obituary for when poets are "retired"
}

type PoetValidationParams struct {
	Designer       string
	SupportedLangs map[string]bool
}

func (p *Poet) Validate(action string, params ...PoetValidationParams) error {
	var (
		err error
	)

	// make sure id, if not an empty string, is a uuid
	if !utils.IsValidUUIDV4(p.Id) && p.Id != "" {
		return fmt.Errorf("Poet Id must be a valid uuid, given %s", p.Id)
	}

	// TODO ensure that only the user namking the create and delete request can perform
	// those actions!
	switch action {
	case consts.CREATE:
		if len(params) == 0 {
			return fmt.Errorf(
				"validation parameters must be provided for %s action",
				consts.CREATE,
			)
		}

		err = p.CheckRequiredFields(params[0])
		if err != nil {
			return err
		}
	case consts.READ:
		// the id *must* be populated and valid
		if p.Id == "" {
			return fmt.Errorf("poet id *must* be provided on %s", consts.READ)
		}
	case consts.UPDATE:
		if len(params) == 0 {
			return fmt.Errorf(
				"validation parameters must be provided for %s action",
				consts.UPDATE,
			)
		}

		// designer must be provided AND match the given validation parameter
		if p.Designer == "" || p.Designer != params[0].Designer {
			return fmt.Errorf("Invalid poet designer provided")
		}

		// the id *must* be populated and valid
		if p.Id == "" {
			return fmt.Errorf("poet id *must* be provided on %s", consts.READ)
		}
	case consts.DELETE:
		if len(params) == 0 {
			return fmt.Errorf(
				"validation parameters must be provided for %s action",
				consts.DELETE,
			)
		}

		// designer must be provided AND match the given validation parameter
		if p.Designer == "" || p.Designer != params[0].Designer {
			return fmt.Errorf("Invalid poet designer provided")
		}

		// the id *must* be populated and valid
		if p.Id == "" {
			return fmt.Errorf("poet id *must* be provided on %s", consts.READ)
		}
	}

	return nil
}

// check required fields for creation
func (p *Poet) CheckRequiredFields(params PoetValidationParams) error {
	var (
		err error
	)

	// we already know that the Id field is valid

	// designer must be provided AND match the given validation parameter
	if p.Designer == "" || p.Designer != params.Designer {
		return fmt.Errorf("Invalid poet designer provided")
	}

	// ensure name is non-empty and obeys naming rules
	err = utils.ValidateUsername(p.Name)
	if err != nil {
		return err
	}

	// limit the size of the description
	if utf8.RuneCountInString(p.Description) > POET_DESCRIPTION_MAX_CHARS {
		return fmt.Errorf("poet description must be below 2k characters")
	}

	// ensure that language is provided and within supported languages
	if _, isSupported := params.SupportedLangs[p.Language]; !isSupported {
		return fmt.Errorf("poet language (%s) not supported (╥﹏╥)", p.Language)
	}

	return nil
}

/*
   db methods
*/

// package level globals for storing prepared sql statements
var (
	poetCreateStmt  *sql.Stmt
	poetReadStmt    *sql.Stmt
	poetReadAllStmt *sql.Stmt
	poetDeleteStmt  *sql.Stmt
)

func CreatePoetsTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS poets (
		          id UUID NOT NULL UNIQUE,
                          designer UUID REFERENCES users NOT NULL,
                          birthDate TIMESTAMP WITH TIME ZONE NOT NULL,
                          deathDate TIMESTAMP WITH TIME ZONE NOT NULL,
                          name VARCHAR(255) NOT NULL UNIQUE,
                          description TEXT NOT NULL,
                          language VARCHAR(255) NOT NULL,
                          execPath VARCHAR(255) NOT NULL UNIQUE,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}

// NOTE [cw|am] 2.21.2018 do we *really* need to be passing in the ID here?
// why can't we just set it in the struct before the function is called??
// that way, we have a cleaner function signature but also have the ability of
// deterministicaly being able to control the value of the ID from outside of
// the function for the sake of testing.
func (p *Poet) Create(id string, db *sql.DB) error {
	var (
		err error
	)

	// assign id
	p.Id = id

	// set birthday
	p.BirthDate = time.Now().Truncate(time.Millisecond)

	// prepare statement if not already done so.
	if poetCreateStmt == nil {
		// create statement
		stmt := `INSERT INTO poets (
                           id, designer, name, birthDate, deathDate, description, language, execPath
                         ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		poetCreateStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	_, err = poetCreateStmt.Exec(
		p.Id,
		p.Designer,
		p.Name,
		p.BirthDate,
		p.DeathDate,
		p.Description,
		p.Language,
		p.ExecPath,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Poet) Read(db *sql.DB) error {
	var (
		err error
	)

	// prepare statement if not already done so.
	if poetReadStmt == nil {
		// read statement
		stmt := `SELECT id, designer, name, birthDate, deathDate, description, language, execPath
                         FROM poets WHERE id = $1`
		poetReadStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	// make sure user Id is actually populated

	// run prepared query over arguments
	err = poetReadStmt.
		QueryRow(p.Id).
		Scan(
			&p.Id,
			&p.Designer,
			&p.Name,
			&p.BirthDate,
			&p.DeathDate,
			&p.Description,
			&p.Language,
			&p.ExecPath,
		)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("No poet with id %s", p.Id)
	case err != nil:
		return err
	}

	return nil
}

// delete should keep meta about poets in the system along with their poems, but
// should remove all files from the file system and assign a death date.
func (p *Poet) Delete(db *sql.DB) error {

	return nil
}

func ReadPoets(db *sql.DB) ([]*Poet, error) {
	var (
		poets []*Poet = []*Poet{}
		err   error
	)

	// prepare statement if not already done so.
	if poetReadAllStmt == nil {
		// readAll statement
		// TODO pagination
		stmt := `SELECT id, designer, name, birthDate, deathDate, description, language, execPath
                         FROM poets`
		poetReadAllStmt, err = db.Prepare(stmt)
		if err != nil {
			return poets, nil
		}
	}

	rows, err := poetReadAllStmt.Query()
	if err != nil {
		return poets, err
	}

	defer rows.Close()
	for rows.Next() {
		poet := &Poet{}
		err = rows.Scan(
			&poet.Id,
			&poet.Designer,
			&poet.Name,
			&poet.BirthDate,
			&poet.DeathDate,
			&poet.Description,
			&poet.Language,
			&poet.ExecPath,
		)
		if err != nil {
			return poets, err
		}

		// append scanned user into list of all users
		poets = append(poets, poet)
	}
	if err := rows.Err(); err != nil {
		return poets, err
	}

	return poets, nil
}
