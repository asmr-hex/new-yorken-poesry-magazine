package types

import (
	"database/sql"
	"time"
)

// Notes about poet executables:
// (1) executables can be named anything, but will be renamed to some standard upon upload
// (2) parameter data can be optionally uploaded if the model decides to store parameters in an external file
// executables will be stored on the filesystem in a safe dir with the path /some/path/bin/<poetId>/

type Poet struct {
	Id          string
	BirthDate   time.Time // so we can show years active
	DeathDate   time.Time // this should be set to null for currently active poets
	Name        string
	Description string
	ExecPath    string // or possibly a Path, this is the path to the source code
	// TODO additional statistics: specifically, it would be cool to see the success rate
	// of a particular poet along with the timeline of how their poems have been recieved
}

func (p *Poet) Validate() error {

	return nil
}

/*
   db methods
*/

func (*Poet) CreateTable(db *sql.DB) error {
	mkTableStmt := `CREATE TABLE IF NOT EXISTS poets (
		          id UUID NOT NULL UNIQUE,
                          birthDate TIMESTAMP WITH TIME ZONE NOT NULL,
                          deathDate TIMESTAMP WITH TIME ZONE,
                          name VARCHAR(255) NOT NULL UNIQUE,
                          description TEXT,
                          execPath VARCHAR(255) NOT NULL UNIQUE,
		          PRIMARY KEY (id)
	)`

	_, err := db.Exec(mkTableStmt)
	if err != nil {
		return err
	}

	return nil
}
