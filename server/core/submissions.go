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
	ExecContext             *env.ExecContext
	NumPoemsInUpcomingIssue int
	Period                  time.Duration
	numConcurrentExecs      int
	wait                    chan bool
	limboDir                string
	db                      *sql.DB
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
	var (
		bestPoems = []*types.Poem{}
		err       error
	)

	issue, err := types.GetUpcomingIssue(s.db)
	if err != nil {
		s.Error(err.Error())
	}

	// for each poem-committe-member pair, generate score

	// read in all limbo files
	files, err := filepath.Glob(filepath.Join(s.limboDir, "*.txt"))
	if err != nil {
		s.Error(err.Error())
	}

	// iterate through filenames in limbo dir
	for _, file := range files {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			s.Error(err.Error())
		}

		// initialize the candidate poem and its scores
		poem := types.Poem{}
		scores := []float64{}

		// since we are running the critics in a go routine pool,
		// we need to create a channel for them to submit their
		// scores to and a channel to tell the score collector
		// when to stop.
		submit := make(chan float64)
		done := make(chan bool)

		err = json.Unmarshal(bytes, &poem)
		if err != nil {
			s.Error(err.Error())
		}

		// define a function that will take a poet judge
		// and will critique the poem in question
		f := func(critic *types.Poet) {
			// when the critic is done, get out of line
			defer func() { <-s.wait }()

			// set execution config
			critic.ExecContext = s.ExecContext

			score, err := critic.CritiquePoem(poem.Content)
			if err != nil {
				s.Error(err.Error())
			}

			// submit this score
			submit <- score
		}

		// start score collector
		go func() {
			for {
				select {
				case score := <-submit:
					// a score has been submitted for this poem
					// add it to the array of scores! Wahooo
					scores = append(scores, score)
				case <-done:
					// all the critics are finished scoring!
					return
				}
			}
		}()

		// ( ͡° ͜ʖ ( ͡° ͜ʖ ( ͡° ͜ʖ ( ͡° ͜ʖ ͡°) ͜ʖ ͡°)ʖ ͡°)ʖ ͡°)
		// send each critic into the critical void (e.g. go routine pool)
		// to critique these poetic verses! ( ･_･)♡
		for _, critic := range issue.Committee {
			// wait in line for judgment
			s.wait <- true

			// execute poet
			go f(critic)
		}

		// the critics have finished scoring!
		done <- true

		// okay, whew (｡-_-｡), now we must average the scores for this poem
		var total float64 = 0
		for _, score := range scores {
			total += score
		}
		poem.Score = total / float64(len(scores))

		// and decide whether it should be inserted into the best poems!
		bestPoems = s.updateBestPoems(bestPoems, &poem)
	}

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

func (s *SubmissionService) updateBestPoems(bestPoems []*types.Poem, poem *types.Poem) []*types.Poem {
	var (
		newBestPoems = []*types.Poem{}
	)

	// if there are enough slots, just add this poem!
	if len(bestPoems) < s.NumPoemsInUpcomingIssue {
		// insert this poem s.t. array is sorted
		for idx, bestPoem := range bestPoems {
			if poem.Score >= bestPoem.Score {
				// this is really easy to parse, right?
				newBestPoems = append(
					newBestPoems,
					append(
						[]*types.Poem{poem},
						newBestPoems[idx:]...,
					)...,
				)

				return newBestPoems
			}

			// otherwise, just keep appending to the newBestPoems
			newBestPoems = append(newBestPoems, bestPoem)
		}
	}

	for _, bestPoem := range bestPoems {
		switch {
		case poem.Score > bestPoem.Score:
			// okay we *must* include this poem in the next issue,
			// it is simply fabulous!
		}
	}

	return newBestPoems
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
