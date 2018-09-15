package types

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	"github.com/frenata/xaqt"
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
	Id                    string           `json:"id"`
	Designer              string           `json:"designer"`            // the writer of the poet (user)
	BirthDate             time.Time        `json:"birthDate"`           // so we can show years active
	DeathDate             time.Time        `json:"deathDate,omitempty"` // this should be set to null for currently active poets
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	Language              string           `json:"language"`
	ProgramFileName       string           `json:"programFileName"`
	ParameterFileName     string           `json:"parameterFileName"`
	ParameterFileIncluded bool             `json:"parameterFileIncluded"`
	Path                  string           `json:"-"` // this is the path to the source code
	ExecContext           *env.ExecContext // inherit from platform config
	Likes                 int              // number of likes this poet has
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

// sanitize out sensitive data from poet
//
func (p *Poet) Sanitize() {
	// for now do nothing...
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
	// Note (cw|4.27.2018) we are using the #ValidateUsername function for
	// validating a poet's Name (not a username)
	// TODO (cw|4.27.2018) proposal to rename (or make another) validation
	// utility function for just Names in general to avoid confusion.
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
                          programFileName VARCHAR(255) NOT NULL,
                          parameterFileName VARCHAR(255) NOT NULL,
                          parameterFileIncluded BOOL NOT NULL,
                          path VARCHAR(255) NOT NULL UNIQUE,
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
                           id, designer, name, birthDate, deathDate, description, language, programFileName, parameterFileName, parameterFileIncluded, path
                         ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
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
		p.ProgramFileName,
		p.ParameterFileName,
		p.ParameterFileIncluded,
		p.Path,
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
		stmt := `SELECT id, designer, name, birthDate, deathDate, description, language, programFileName, parameterFileName, parameterFileIncluded, path
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
			&p.ProgramFileName,
			&p.ParameterFileName,
			&p.ParameterFileIncluded,
			&p.Path,
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

func CountPoets(db *sql.DB) (int, error) {
	var (
		count int
		err   error
	)

	err = db.QueryRow(`SELECT COUNT(*) FROM poets;`).Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func SelectRandomPoets(n int, db *sql.DB) ([]*Poet, error) {
	var (
		poets []*Poet
		err   error
	)

	rows, err := db.Query(`
            SELECT id, designer, name, birthDate, deathDate, description, language, programFileName, parameterFileName, parameterFileIncluded, path
            FROM poets
            ORDER BY RANDOM()
            LIMIT $1
        `, n)
	if err != nil {
		return nil, err
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
			&poet.ProgramFileName,
			&poet.ParameterFileName,
			&poet.ParameterFileIncluded,
			&poet.Path,
		)
		if err != nil {
			return poets, err
		}

		// append scanned user into list of all poets
		poets = append(poets, poet)
	}
	if err := rows.Err(); err != nil {
		return poets, err
	}

	return poets, nil
}

func GetUnderdogPoets(n int, db *sql.DB) ([]*Poet, error) {
	// TODO (cw|9.14.2018) choose underdogs more equitably:
	// * choose based on least popular programming languages
	// * choose based on how new the poets are
	// * choose based on if they haven't had poems published
	// * choose based on if they have had poems published, but with low scores

	// for now just select randomly
	return SelectRandomPoets(n, db)
}

func GetFancyPoets(n int, db *sql.DB) ([]*Poet, error) {
	// var (
	// 	poets []*Poet
	// 	err   error
	// )

	// select most prolific poets, order by their average quality

	// maxPubs

	// TODO (cw|9.14.2018) use a good metric for this combination score
	// of prolificness and quality!

	// For now, just randomly choose -___- but change that!!!
	return SelectRandomPoets(n, db)
}

func ReadPoets(db *sql.DB, filter ...string) ([]*Poet, error) {
	var (
		poets []*Poet = []*Poet{}
		err   error
	)

	// prepare statement if not already done so.
	if poetReadAllStmt == nil {
		// readAll statement
		// TODO pagination
		stmt := `SELECT id, designer, name, birthDate, deathDate, description, language, programFileName, parameterFileName, parameterFileIncluded, path
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
			&poet.ProgramFileName,
			&poet.ParameterFileName,
			&poet.ParameterFileIncluded,
			&poet.Path,
		)
		if err != nil {
			return poets, err
		}

		// append scanned user into list of all poets
		poets = append(poets, poet)
	}
	if err := rows.Err(); err != nil {
		return poets, err
	}

	return poets, nil
}

func (p *Poet) GeneratePoem() (*Poem, error) {
	var (
		poem *Poem
		err  error
	)

	ctx, code, err := p.setupExecutionSandbox()
	if err != nil {
		return nil, err
	}

	// execute poem generation task
	results, _ := ctx.Evaluate(p.Language, code, []string{"write"})
	// TODO (cw|9.2.2018) Evaluate returns an xaqt.Message which we shoul use
	// to extract the appropriate error.

	poem = &Poem{
		Date:    time.Now(),
		Author:  p,
		Content: results[0],
	}

	return poem, nil
}

func (p *Poet) CritiquePoem(poem string) (float64, error) {
	var (
		score float64
		err   error
	)

	ctx, code, err := p.setupExecutionSandbox()
	if err != nil {
		return score, err
	}

	results, _ := ctx.Evaluate(p.Language, code, []string{"critique"})

	score, err = strconv.ParseFloat(results[0], 64)
	if err != nil {
		return score, err
	}

	// ensure that the score is within the range [0, 1]
	if !(score >= 0 && score <= 1) {
		return score, fmt.Errorf("invalid score (%f) given by poet %s", score, p.Name)
	}

	return score, nil
}

func (p *Poet) StudyPoem(poem string) (bool, error) {
	ctx, code, err := p.setupExecutionSandbox()
	if err != nil {
		return false, err
	}

	results, _ := ctx.Evaluate(p.Language, code, []string{"study"})

	success, err := strconv.ParseBool(results[0])
	if err != nil {
		return false, err
	}

	return success, nil
}

// creates and configures the execution sandbox.
func (p *Poet) setupExecutionSandbox() (*xaqt.Context, xaqt.Code, error) {
	var (
		ctx  *xaqt.Context
		code xaqt.Code
		err  error
	)

	if p.ExecContext == nil {
		return nil, code, fmt.Errorf("developer error! exec context not set for poet %s", p.Id)
	}

	// setup execution context
	ctx, err = xaqt.NewContext(
		xaqt.GetCompilers(),
		xaqt.ExecDir(p.ExecContext.Dir),
		xaqt.ExecMountDir(p.ExecContext.MountDir),
	)
	if err != nil {
		return nil, code, err
	}

	code = xaqt.Code{
		IsFile:            true,
		SourceFileName:    p.ProgramFileName,
		ResourceFileNames: []string{},
		Path:              p.Path,
	}

	if p.ParameterFileIncluded {
		code.ResourceFileNames = []string{p.ParameterFileName}
	}

	return ctx, code, nil
}

func (p *Poet) TestPoet() error {
	// test poem generation
	poem, err := p.GeneratePoem()
	if err != nil {
		return err
	}

	if strings.TrimSpace(poem.Content) == "" {
		return fmt.Errorf("%s's work is simply vapid.", p.Name)
	}

	// test poem evaluation
	score, err := p.CritiquePoem(consts.THE_BEST_POEM)
	if err != nil {
		return err
	}

	if score < 0 || score > 1 {
		return fmt.Errorf("%s is a bad critic.", p.Name)
	}

	// test self updating
	success, err := p.StudyPoem(consts.THE_BEST_POEM)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("%s is unable to learn.", p.Name)
	}

	return nil
}
