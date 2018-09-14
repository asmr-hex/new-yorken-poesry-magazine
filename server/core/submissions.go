package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
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

	// now we have a map of the n highest rated poems in bestPoems

	// store these poems in the database
	for _, poem := range bestPoems {
		poem.Issue = issue
		err = poem.Create(s.db)
		if err != nil {
			s.Error(err.Error())
		}
	}
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

	// optimization! first check last element and if this score is lower,
	// then we just return the original array.
	if bestPoems[len(bestPoems)-1].Score > poem.Score {
		return bestPoems
	}

	// once we find a place to insert this, we need to get a list of the
	// lowest scores in the bestPoems list. We will then randomly shave off
	// one of the entries. fair is fair.
	lowestScoreBlock := []*types.Poem{}
	lowestScore := 0.0

	// TODO (cw|9.13.2018) I can probably use a min-heap for this...
	for _, bestPoem := range bestPoems {
		if len(lowestScoreBlock) == 0 {
			// we still haven't found where this poem should reside in
			// the sorted bestPoems array...
			switch {
			case poem.Score > bestPoem.Score:
				// alright, we know our candidate poem is *at least*
				// better than this bestPoem...
				lowestScore = bestPoem.Score
				lowestScoreBlock = []*types.Poem{bestPoem}
			case poem.Score == bestPoem.Score:
				// alright, we know that the candidate poem is *just*
				// as good as this bestPoem
				lowestScore = poem.Score
				lowestScoreBlock = []*types.Poem{poem, bestPoem}
			default:
				newBestPoems = append(newBestPoems, bestPoem)
			}
		} else {
			// all bestPoems are gauranteed to have lower scores now
			// than poem
			if bestPoem.Score < lowestScore {
				// append out-dated lowestScoreBlock to result
				newBestPoems = append(newBestPoems, lowestScoreBlock...)

				// update lowestScoreBlock and lowestScore
				lowestScoreBlock = []*types.Poem{bestPoem}
				lowestScore = bestPoem.Score
			} else {
				// if we are here, this bestPoem *must* have a score
				// equal to the current lowestScore.
				lowestScoreBlock = append(lowestScoreBlock, bestPoem)
			}
		}
	}

	// okay, at this point we have all the poems sharded into two arrays
	// (1) lowestScoreBlock: all the poems with the same lowest score
	// (2) newBestPoems: the rest of the non-lowestScoreBlock poems (sorted)
	// we need to randomly drop on item from the lowestScoreBlock
	rand.Seed(time.Now().Unix())
	rmIdx := rand.Intn(len(lowestScoreBlock))

	// delete a random element from the lowestScoreBlock
	lowestScoreBlock = append(lowestScoreBlock[:rmIdx], lowestScoreBlock[rmIdx+1:]...)

	// stitch back together the newBestPoems + lowestScoreBlock
	newBestPoems = append(newBestPoems, lowestScoreBlock...)

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
