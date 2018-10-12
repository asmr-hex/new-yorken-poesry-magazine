package types

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/consts"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	"github.com/frenata/xaqt"
	"github.com/lib/pq"
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
	Designer              *User            `json:"designer"`            // the writer of the poet (user)
	BirthDate             time.Time        `json:"birthDate"`           // so we can show years active
	DeathDate             time.Time        `json:"deathDate,omitempty"` // this should be set to null for currently active poets
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	Language              string           `json:"language"`
	ProgramFileName       string           `json:"programFileName"`
	ParameterFileName     string           `json:"parameterFileName"`
	ParameterFileIncluded bool             `json:"-"`
	Path                  string           `json:"-"` // this is the path to the source code
	ExecContext           *env.ExecContext `json:"-"` // inherit from platform config
	Likes                 int              // number of likes this poet has
	// TODO additional statistics: specifically, it would be cool to see the success rate
	// of a particular poet along with the timeline of how their poems have been recieved

	// what if we also had a poet obituary for when poets are "retired"
}

// this struct is strictly for extracting possibly null valued
// fields from the database -___-
// we will only use this struct if we are OUTER JOINING poets on
// another table (e.g. users, since some users might not have poets)
// TODO (cw|9.20.2018) figure out a better way to do this...
type PoetNullable struct {
	Id                    sql.NullString
	DesignerId            sql.NullString
	BirthDate             pq.NullTime
	DeathDate             pq.NullTime
	Name                  sql.NullString
	Description           sql.NullString
	Language              sql.NullString
	ProgramFileName       sql.NullString
	ParameterFileName     sql.NullString
	ParameterFileIncluded sql.NullBool
	Path                  sql.NullString
}

func (pn *PoetNullable) Convert() *Poet {
	return &Poet{
		Id:                    pn.Id.String,
		Designer:              &User{Id: pn.DesignerId.String},
		BirthDate:             pn.BirthDate.Time,
		DeathDate:             pn.DeathDate.Time,
		Name:                  pn.Name.String,
		Description:           pn.Description.String,
		Language:              pn.Language.String,
		Path:                  pn.Path.String,
		ProgramFileName:       pn.ProgramFileName.String,
		ParameterFileName:     pn.ParameterFileName.String,
		ParameterFileIncluded: pn.ParameterFileIncluded.Bool,
	}
}

type PoetValidationParams struct {
	DesignerId     string
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
		if p.Designer == nil || p.Designer.Id != params[0].DesignerId {
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
		if p.Designer == nil || p.Designer.Id != params[0].DesignerId {
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
	if p.Designer == nil || p.Designer.Id != params.DesignerId {
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
	poetCreateStmt     *sql.Stmt
	poetReadStmt       *sql.Stmt
	poetReadAllStmt    *sql.Stmt
	poetDeleteStmt     *sql.Stmt
	poetCodeReadStmt   *sql.Stmt
	countPoetPoemsStmt *sql.Stmt
	deletePoetStmt     *sql.Stmt
	softDeletePoetStmt *sql.Stmt
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
                          deleted BOOL NOT NULL DEFAULT false,
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
func (p *Poet) Create(db *sql.DB) error {
	var (
		err error
	)

	// assume id is already assigned

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
		p.Designer.Id,
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
		stmt := `
                         SELECT p.id, name, birthDate, deathDate, description,
                                language, programFileName, parameterFileName,
                                parameterFileIncluded, path,
                                u.id, u.username, u.email
                         FROM poets p
                         INNER JOIN users u ON (p.designer = u.id)
                         WHERE p.id = $1
                `
		poetReadStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	// initialize designer struct before accessing it
	p.Designer = &User{}

	// run prepared query over arguments
	err = poetReadStmt.
		QueryRow(p.Id).
		Scan(
			&p.Id,
			&p.Name,
			&p.BirthDate,
			&p.DeathDate,
			&p.Description,
			&p.Language,
			&p.ProgramFileName,
			&p.ParameterFileName,
			&p.ParameterFileIncluded,
			&p.Path,
			&p.Designer.Id,
			&p.Designer.Username,
			&p.Designer.Email,
		)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("No poet with id %s", p.Id)
	case err != nil:
		return err
	}

	return nil
}

func (p *Poet) ReadCode(db *sql.DB) (*Code, error) {
	var (
		code *Code = &Code{}
		err  error
	)

	// prepare statement if not already done so.
	if poetCodeReadStmt == nil {
		// read statement
		stmt := `
                         SELECT id, language, programFileName, path
                         FROM poets
                         WHERE id = $1 AND deleted = false
                `
		poetCodeReadStmt, err = db.Prepare(stmt)
		if err != nil {
			return nil, err
		}
	}

	// run prepared query over arguments
	err = poetCodeReadStmt.
		QueryRow(p.Id).
		Scan(
			&code.Author,
			&code.Language,
			&code.FileName,
			&code.Path,
		)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("No poet with id %s", p.Id)
	case err != nil:
		return nil, err
	}

	// read code in from file
	err = code.Read()
	if err != nil {
		return nil, err
	}

	return code, nil
}

// hard delete if a poet has not published anything. soft-delete if a poet has published...
func (p *Poet) Delete(db *sql.DB) error {
	var (
		poems int
		err   error
	)

	// assume that we delete the poet directory within the caller scope

	if countPoetPoemsStmt == nil {
		stmt := `SELECT COUNT(*) FROM poems p WHERE p.author = $1`
		countPoetPoemsStmt, err = db.Prepare(stmt)
		if err != nil {
			return err
		}
	}

	// how many poems has this poet published?
	err = countPoetPoemsStmt.QueryRow(p.Id).Scan(&poems)
	if err != nil {
		return err
	}

	// if the poet has published at least 1 poem, we will keep its meta-data up here.
	// otherwise, we will delete the entire record.
	if poems == 0 {
		// delete this poet and its files
		if deletePoetStmt == nil {
			stmt := `DELETE FROM poets WHERE id = $1`
			deletePoetStmt, err = db.Prepare(stmt)
			if err != nil {
				return err
			}
		}

		_, err = deletePoetStmt.Exec(p.Id)
		if err != nil {
			return err
		}

	} else {
		// delete the files and update the deathdate
		if softDeletePoetStmt == nil {
			stmt := `UPDATE poets SET deleted = true, deathDate = $2  WHERE id = $1`
			softDeletePoetStmt, err = db.Prepare(stmt)
			if err != nil {
				return err
			}
		}

		_, err = softDeletePoetStmt.Exec(p.Id, time.Now())
		if err != nil {
			return err
		}

	}

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

func SelectRandomPoets(n int, db *sql.DB, exclude ...*Poet) ([]*Poet, error) {
	var (
		poets []*Poet
		err   error
	)

	setString := ""
	whereClause := ""
	for idx, poet := range exclude {
		if idx == 0 {
			setString = fmt.Sprintf(`'%s'`, poet.Id)
		} else {
			setString += fmt.Sprintf(`, '%s'`, poet.Id)
		}
	}
	if setString != "" {
		whereClause = fmt.Sprintf("WHERE p.id NOT IN (%s)", setString)
	}
	sql := fmt.Sprintf(`
                         SELECT p.id, name, birthDate, deathDate, description,
                                language, programFileName, parameterFileName,
                                parameterFileIncluded, path,
                                u.id, u.username, u.email
                         FROM poets p
                         INNER JOIN users u ON (p.designer = u.id)
                         %s
                         ORDER BY RANDOM()
                         LIMIT $1
        `, whereClause)

	rows, err := db.Query(sql, n)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		poet := &Poet{Designer: &User{}}
		err = rows.Scan(
			&poet.Id,
			&poet.Name,
			&poet.BirthDate,
			&poet.DeathDate,
			&poet.Description,
			&poet.Language,
			&poet.ProgramFileName,
			&poet.ParameterFileName,
			&poet.ParameterFileIncluded,
			&poet.Path,
			&poet.Designer.Id,
			&poet.Designer.Username,
			&poet.Designer.Email,
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

func GetUnderdogPoets(n int, db *sql.DB, exclude ...*Poet) ([]*Poet, error) {
	// TODO (cw|9.14.2018) choose underdogs more equitably:
	// * choose based on least popular programming languages
	// * choose based on how new the poets are
	// * choose based on if they haven't had poems published
	// * choose based on if they have had poems published, but with low scores

	// for now just select randomly
	return SelectRandomPoets(n, db, exclude...)
}

func GetFancyPoets(n int, db *sql.DB, exclude ...*Poet) ([]*Poet, error) {
	// var (
	// 	poets []*Poet
	// 	err   error
	// )

	// select most prolific poets, order by their average quality

	// maxPubs

	// TODO (cw|9.14.2018) use a good metric for this combination score
	// of prolificness and quality!

	// For now, just randomly choose -___- but change that!!!
	return SelectRandomPoets(n, db, exclude...)
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
		stmt := `
                         SELECT p.id, name, birthDate, deathDate, description,
                                language, programFileName, parameterFileName,
                                parameterFileIncluded, path,
                                u.id, u.username, u.email
                         FROM poets p
                         INNER JOIN users u ON (p.designer = u.id)
                `
		poetReadAllStmt, err = db.Prepare(stmt)
		if err != nil {
			return nil, err
		}
	}

	rows, err := poetReadAllStmt.Query()
	if err != nil {
		return poets, err
	}

	defer rows.Close()
	for rows.Next() {
		poet := &Poet{Designer: &User{}}
		err = rows.Scan(
			&poet.Id,
			&poet.Name,
			&poet.BirthDate,
			&poet.DeathDate,
			&poet.Description,
			&poet.Language,
			&poet.ProgramFileName,
			&poet.ParameterFileName,
			&poet.ParameterFileIncluded,
			&poet.Path,
			&poet.Designer.Id,
			&poet.Designer.Username,
			&poet.Designer.Email,
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
	results, msg := ctx.Evaluate(p.Language, code, PoetAPIGenerateArgs())
	// TODO (cw|9.2.2018) Evaluate returns an xaqt.Message which we shoul use
	// to extract the appropriate error.
	// NOTE (cw|10.12.2018) this is the message that can be used to populate the
	// proposed Poet Error Console described in issue #54.
	fmt.Println(msg)

	rawPoem, err := p.ParseRawPoem(results)
	if err != nil {
		return nil, err
	}

	poem = &Poem{
		Date:    time.Now(),
		Author:  p,
		Title:   rawPoem.Title,
		Content: rawPoem.Content,
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

	results, _ := ctx.Evaluate(p.Language, code, PoetAPICritiqueArgs(poem))

	critique, err := p.ParseRawCritique(results)
	if err != nil {
		return score, err
	}

	return critique.Score, nil
}

func (p *Poet) StudyPoems(poems ...string) (bool, error) {
	ctx, code, err := p.setupExecutionSandbox()
	if err != nil {
		return false, err
	}

	results, _ := ctx.Evaluate(p.Language, code, PoetAPIUpdateArgs(poems...))

	rawUpdate, err := p.ParseRawUpdate(results)
	if err != nil {
		return false, err
	}

	return rawUpdate.Success, nil
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
		xaqt.InputType(xaqt.ArgsInput),
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
	success, err := p.StudyPoems(consts.THE_BEST_POEM)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("%s is unable to learn.", p.Name)
	}

	return nil
}
