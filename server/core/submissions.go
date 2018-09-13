package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/env"
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
)

const (
	NUM_CONCURRENT_EXECS_DEFAULT = 5
)

type SubmissionService struct {
	*Logger
	ExecContext        *env.ExecContext
	Period             time.Duration
	numConcurrentExecs int
	wait               chan bool
	limboDir           string
	db                 *sql.DB
}

func NewSubmissions(execCtx *env.ExecContext, db *sql.DB) *SubmissionService {
	return &SubmissionService{
		Logger:             NewLogger(os.Stdout),
		ExecContext:        execCtx,
		Period:             time.Hour * 24 * 7,
		numConcurrentExecs: NUM_CONCURRENT_EXECS_DEFAULT,
		wait:               make(chan bool, NUM_CONCURRENT_EXECS_DEFAULT),
		limboDir:           "/poems-awaiting-judgment",
		db:                 db,
	}
}

func (s *SubmissionService) UpdatePeriod(period time.Duration) {
	s.Period = period
}

func (s *SubmissionService) UpdateNumberOfConcurrentExecs(n int) {
	s.numConcurrentExecs = n
	newWait := make(chan bool, n)

	// TODO (cw|9.12.2018) need to get a write lock before draining...

	// drain old channel
	for len(s.wait) > 0 {
		<-s.wait
		newWait <- true
	}

	// assign new wait
	s.wait = newWait
}

func (s *SubmissionService) StartScheduler() {
	ticker := time.NewTicker(s.Period)

	// that's right, we plan on publishing this magazine ~*~*~*F O R E V E R*~*~*~
	// just look at this unqualified loop! It's like staring into the void of perpetual
	// poetical motion, like an unbreakable möbius band, lithe yet oppressive.
	for {
		<-ticker.C

		s.OpenCallForSubmissions()

		s.SelectWinningPoems()

		s.ReleaseNewIssue()

		s.ChooseNewCommitteeMembers()

		s.CleanUp()

		s.AllowPoetsToLearn()
	}
}

func (s *SubmissionService) OpenCallForSubmissions() {
	var (
		poets []*types.Poet
		err   error
	)

	// get all candidate poets TODO (cw|9.12.2018) filter this read for only active poets
	poets, err = types.ReadPoets(s.db)
	if err != nil {
		// this may happen if the database goes down or something...
		// but it shouldn't normally happen  ¯\_(ツ)_/¯
		s.Error(err.Error())
	}

	// run each candidate poet
	s.ElicitPoemsFrom(poets)
}

func (s *SubmissionService) ElicitPoemsFrom(poets []*types.Poet) {

	f := func(poet *types.Poet) error {
		// when the poet is complete, get out of line
		defer func() { <-s.wait }()

		// set execution config
		poet.ExecContext = s.ExecContext

		// ask this poet to write some verses
		poem, err := poet.GeneratePoem()
		if err != nil {
			s.Error(err.Error())
			return err
		}

		// store this poem in the limbo directory so it can be judged

		// Note: we are going through the trouble of marshalling into json
		// and storing on the fs because we don't really care about how long
		// this takes...i mean, give us a break, CFPs for humans are long
		// and arduous administrative processes. why do *we* also have to be
		// optimally performant???
		filename := filepath.Join(s.limboDir, fmt.Sprintf("%s.txt", poet.Id))
		bytes, err := json.Marshal(poem)
		if err != nil {
			s.Error(err.Error())
			return err
		}

		err = ioutil.WriteFile(filename, bytes, 0700)
		if err != nil {
			s.Error(err.Error())
			return err
		}

		return nil
	}

	for _, poet := range poets {
		// wait in line for execution
		s.wait <- true

		// execute poet
		go f(poet)
	}
}

func (s *SubmissionService) SelectWinningPoems() {
	issue, err := types.GetUpcomingIssue(s.db)
	if err != nil {
		s.Error(err.Error())
	}

	_ = issue

	// for each poem-committe-member pair, generate score

	// iterate through filenames in limbo dir
	//   - unmarshal each into a Poem struct
	//   - have each committee member judge this poem
	//   - average each committee members scores
	//   - conditionally add this poem to a global map of top-ranked poems
	//     - does this poem score fall within the range of high/low scores in map?
	//     - or is there still room in the map?
	//       - place poem in map and kick out lower scored poem if necessary
	//         - if poem has same score as k lowest poems, flip a coin to pick which one to kick out (including the new candidate)

	// now we have a map of the n highest rated poems

	// store these poems in the database
}

func (s *SubmissionService) ReleaseNewIssue() {
	// release new issue
}

func (s *SubmissionService) ChooseNewCommitteeMembers() {
	// choose new committee members
}

func (s *SubmissionService) AllowPoetsToLearn() {
	// do stuff
}

func (s *SubmissionService) CleanUp() {
	// TODO (cw|9.12.2018) clear limboDir
}
