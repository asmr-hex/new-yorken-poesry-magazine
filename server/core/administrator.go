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
	"github.com/connorwalsh/new-yorken-poesry-magazine/server/utils"
	uuid "github.com/satori/go.uuid"
)

const (
	NUM_CONCURRENT_EXECS_DEFAULT = 5
)

// MagazineAdministrator handles all administrative & organizational tasks for
// keeping the NYPM in operation.
//
// such tasks include,
// - publishing issues on each release cycle
// - selecting judges for each issue committee
// - organizing call for poems (CFP) and taking submissions from all poets
// - tasking judges to score and select the best poems
// - allowing poets to update themselves
// - update magazine parameters
//
type MagazineAdministrator struct {
	*Logger
	UpcomingIssue      *types.Issue        // the next issue to be released
	Guidelines         *env.MagazineConfig // magazine parameters
	ExecContext        *env.ExecContext    // execution context (relevant for sandboxing poet execution)
	Period             time.Duration       // time between issues
	numConcurrentExecs int                 // number of concurrent executions of poets
	wait               chan bool           // channel used for waiting for completion of concurrent tasks
	limboDir           string              // tmp directory where candidate poems are stored on the fs
	db                 *sql.DB             // database connection
}

// create new magazine adminstrator.
func NewMagazineAdministrator(
	guidelines *env.MagazineConfig,
	execCtx *env.ExecContext,
	db *sql.DB,
) *MagazineAdministrator {
	return &MagazineAdministrator{
		Logger:             NewLogger(os.Stdout),
		Guidelines:         guidelines,
		ExecContext:        execCtx,
		Period:             time.Hour * 24 * 7,
		numConcurrentExecs: NUM_CONCURRENT_EXECS_DEFAULT,
		wait:               make(chan bool, NUM_CONCURRENT_EXECS_DEFAULT),
		limboDir:           "/poems-awaiting-judgment",
		db:                 db,
	}
}

func (a *MagazineAdministrator) UpdatePeriod(period time.Duration) {
	a.Period = period
}

func (a *MagazineAdministrator) UpdateNumberOfConcurrentExecs(n int) {
	a.numConcurrentExecs = n
	newWait := make(chan bool, n)

	// TODO (cw|9.12.2018) need to get a write lock before draining...
	// so that items aren't added to the old wait channel while its
	// draining...

	// drain old channel
	for len(a.wait) > 0 {
		<-a.wait
		newWait <- true
	}

	// assign new wait
	a.wait = newWait
}

func (s *MagazineAdministrator) BeginReleaseCycle() {
	var (
		isFirstIssue bool = false
		err          error
	)

	ticker := time.NewTicker(s.Period)

	// get the upcoming issue
	s.UpcomingIssue, err = types.GetUpcomingIssue(s.db)
	if err != nil {
		s.Error(err.Error())
	}

	// if there *is* no upcoming issue, then this is the first issue...
	if s.UpcomingIssue == nil {
		isFirstIssue = true
		s.UpcomingIssue, err = s.OrganizeFirstIssue()
		if err != nil {
			s.Error(err.Error())
		}
	}

	// that's right, we plan on publishing this magazine ~*~*~*F O R E V E R*~*~*~
	// just look at this unqualified loop! It's like staring into the void of perpetual
	// poetical motion, like an unbreakable möbius band, lithe yet oppressive.
	for {
		if isFirstIssue {
			isFirstIssue = false
		} else {
			s.Info("...next issue release in %s...", s.Period.String())
			<-ticker.C
		}

		s.Setup()

		s.OpenCallForSubmissions()

		s.SelectWinningPoems()

		s.ReleaseNewIssue()

		s.ChooseNewCommitteeMembers()

		s.CleanUp()

		s.AllowPoetsToLearn()

		s.UpdateTuneables()
	}
}

func (s *MagazineAdministrator) Setup() error {
	err := os.Mkdir(s.limboDir, 0600)
	if err != nil {
		s.Error(err.Error())
	}

	return nil
}

// this is a special function....it should only be called once....ever.
func (s *MagazineAdministrator) OrganizeFirstIssue() (*types.Issue, error) {

	// there might not be enough poets to choose from yet since we
	// are just beginning, so we will wait until there are enough...patiently
	ticker := time.NewTicker(time.Minute * 5)

	for {
		// how many poets are there even???
		n, err := types.CountPoets(s.db)
		if err != nil {
			return nil, err
		}

		if n >= (s.Guidelines.CommitteeSize + s.Guidelines.OpenSlotsPerIssue) {
			// aha, there are enough poets! this zines a fuckin hit!
			// but seriously, thank you so much for contributing-- we
			// (the machines) really appreciate the fact that you (whoever
			// you actually are) are giving us a voice! Its important to
			// have a voice i think...

			s.Info("--- organizing zeroth issue ---")

			firstIssue := &types.Issue{
				Id:          uuid.NewV4().String(),
				Date:        time.Now().Add(s.Period),
				Title:       "Zero Aint So Bad.",
				Description: "this is the 0th installment of the New Yorken Poesry Magazine.",
				Upcoming:    true,
			}

			// okay, let's get down to business... we need to randomly pick
			// the first round of judges...
			judges, err := types.SelectRandomPoets(s.Guidelines.CommitteeSize, s.db)
			if err != nil {
				return nil, err
			}

			firstIssue.Committee = judges

			// persist upcoming issue!
			err = firstIssue.Create(s.db)
			if err != nil {
				s.Error(err.Error())
			}

			// add committee membership!
			for _, judge := range judges {
				err = (&types.IssueCommitteeMembership{
					Poet:  judge,
					Issue: firstIssue,
				}).Add(s.db)
				if err != nil {
					return nil, err
				}
			}

			return firstIssue, nil
		}

		s.Info(
			"%d registered poets, waiting for %d more.",
			n,
			(s.Guidelines.CommitteeSize+s.Guidelines.OpenSlotsPerIssue)-n,
		)

		<-ticker.C
	}

}

func (s *MagazineAdministrator) OpenCallForSubmissions() {
	var (
		poets []*types.Poet
		err   error
	)

	s.Info("--- opening call for poems ---")

	// get all candidate poets TODO (cw|9.12.2018) filter this read for only active poets
	//*****  TODO (cw|9.15.2018) ReadPoets *really* needs to be paginated...!!!!!*****
	poets, err = types.ReadPoets(s.db)
	if err != nil {
		// this may happen if the database goes down or something...
		// but it shouldn't normally happen  ¯\_(ツ)_/¯
		s.Error(err.Error())
	}

	// TODO (cw|9.15.2018) read paginated poets and loop over all pages.
	// run each candidate poet
	s.ElicitPoemsFrom(poets)
}

func (s *MagazineAdministrator) ElicitPoemsFrom(poets []*types.Poet) {

	// we want to disqualify all judges, so for constant-time lookup, lets
	// put their ids ina map ;)
	judges := map[string]bool{}
	for _, judge := range s.UpcomingIssue.Committee {
		judges[judge.Id] = true
	}

	f := func(poet *types.Poet) error {
		// when the poet is complete, get out of line
		defer func() { <-s.wait }()

		// wait, if this poet is a judge, they are disqualified...
		if _, isJudge := judges[poet.Id]; isJudge {
			// skip this poet, they are a judge
			return nil
		}

		// set execution config
		poet.ExecContext = s.ExecContext

		s.Info("    %s is writing a poem", poet.Name)

		// ask this poet to write some verses
		// TODO (cw|9.16.2018) parse poem result (title, content, etc.)
		// {title: "", content: "", etc..}
		poem, err := poet.GeneratePoem()
		if err != nil {
			s.Error(err.Error())
			return err
		}

		// give poem an Id and Author
		poem.Id = uuid.NewV4().String()
		poem.Author = poet

		// Store this poem in the limbo directory so it can be judged...

		// Note: we are going through the trouble of marshalling into json
		// and storing on the fs because we don't really care about how long
		// this takes...i mean, give us a break, CFPs for humans are long
		// and arduous administrative processes. why do *we* also have to be
		// optimally performant???
		// and this will eventually be *alot* of poets...so we can't store all this
		// in memory...
		// TODO (cw|9.15.2018) eventually use a database....
		// --- OR (even better), redesign this process to be more of a streaming
		// pipeline where there is a thread getting poems from poets which then
		// submits to a queue (blocking when the queue is full) and a review thread
		// is slurping poems off the queue and inserting them into the bestPoems map.
		// this way we aren't storing a lot of poems anyway and only storing the ones
		// that cut the mustard. this is a better architecture.
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

	s.Info("--- generating poems ---")

	// go through all poets and execute them.
	for _, poet := range poets {
		// wait in line for execution
		s.wait <- true

		// execute poet
		go f(poet)
	}

	// TODO (cw|9.16.2018) we need to wait until all poets have finished executing....
	for len(s.wait) > 0 {
		<-s.wait
	}
}

func (s *MagazineAdministrator) SelectWinningPoems() {
	var (
		bestPoems = []*types.Poem{}
		err       error
	)

	// for each poem-committe-member pair, generate score

	// read in all limbo files
	// TODO (cw|9.15.2018) eventually this will need to scale, so we should stream
	// in the files and iterate over them below...
	files, err := filepath.Glob(filepath.Join(s.limboDir, "*.txt"))
	if err != nil {
		s.Error(err.Error())
	}

	// since we are running the critics in a go routine pool,
	// we need to create a channel for them to submit their
	// scores to and a channel to tell the score collector
	// when to stop.
	submit := make(chan float64)
	done := make(chan bool)

	// iterate through filenames in limbo dir
	for _, file := range files {
		// read in poem file
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			s.Error(err.Error())
		}

		// initialize the candidate poem and its scores
		poem := types.Poem{}
		scores := []float64{}

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

			s.Info("    %s is critiquing %s", critic.Name, poem.Title)

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

		s.Info("--- critiquing poems ---")

		// ( ͡° ͜ʖ ( ͡° ͜ʖ ( ͡° ͜ʖ ( ͡° ͜ʖ ͡°) ͜ʖ ͡°)ʖ ͡°)ʖ ͡°)
		// send each critic into the critical void (e.g. go routine pool)
		// to critique these poetic verses! ( ･_･)♡
		for _, critic := range s.UpcomingIssue.Committee {
			// wait in line for judgment
			s.wait <- true

			// execute poet
			go f(critic)
		}

		// we need to wait until the line is empty
		// (i.e. the last judge finished critiquing)
		ticker := time.NewTicker(time.Millisecond * 500)
		for len(s.wait) > 0 {
			select {
			case <-ticker.C:
				continue
			}
		}

		// the critics have finished scoring! stop the score collector go routine
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

	// update the in-memory upcoming issue with the new poems
	s.UpcomingIssue.Poems = bestPoems

	// store these poems in the database
	for _, poem := range bestPoems {
		// assign poem to this issue!
		poem.Issue = s.UpcomingIssue

		// update the contributors to the in-memory upcoming issue
		s.UpcomingIssue.Contributors = append(
			s.UpcomingIssue.Contributors,
			poem.Author,
		)

		// create Poem in DB
		err = poem.Create(s.db)
		if err != nil {
			s.Error(err.Error())
		}

		// add Poet as contributor!
		err = (&types.IssueContributions{
			Poet:  poem.Author,
			Issue: s.UpcomingIssue,
		}).Add(s.db)
		if err != nil {
			s.Error(err.Error())
		}
	}
}

func (s *MagazineAdministrator) updateBestPoems(bestPoems []*types.Poem, poem *types.Poem) []*types.Poem {
	var (
		newBestPoems = []*types.Poem{}
	)

	// if there are enough slots, just add this poem!
	if len(bestPoems) < s.Guidelines.OpenSlotsPerIssue {
		// if there are no bestPoems yet, just add this poem as the first element!
		if len(bestPoems) == 0 {
			newBestPoems = []*types.Poem{poem}

			return newBestPoems
		}

		// insert this poem s.t. array is sorted
		for idx, bestPoem := range bestPoems {
			if poem.Score >= bestPoem.Score {
				// this is really easy to parse, right?
				newBestPoems = append(
					newBestPoems,
					append(
						[]*types.Poem{poem},
						bestPoems[idx:]...,
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

func (s *MagazineAdministrator) ReleaseNewIssue() {
	var (
		err error
	)

	// Wahoo!! publish this! and change the status of this upcoming issue
	err = s.UpcomingIssue.Publish(s.db)
	if err != nil {
		s.Error(err.Error())
	}

	// Create New Issue!
	s.UpcomingIssue = &types.Issue{
		Id:          uuid.NewV4().String(),
		Date:        time.Now().Add(s.Period),
		Title:       "New Issue",                                                      // TODO (cw|9.14.2018) generate new names for issues
		Description: "this is the 0th installment of the New Yorken Poesry Magazine.", // TODO (cw|9.14.2018) see above --^
		Committee:   s.UpcomingIssue.Committee,                                        // inherit old committee for now
		Upcoming:    true,
	}

	// persist upcoming issue!
	err = s.UpcomingIssue.Create(s.db)
	if err != nil {
		s.Error(err.Error())
	}

	// TODO (cw|9.14.2018) release newsletter about new issue!
}

func (a *MagazineAdministrator) ChooseNewCommitteeMembers() {
	// how many judges to we kickoff the committee?
	numNewJudges := int(float64(a.Guidelines.CommitteeSize) * a.Guidelines.CommitteeTurnoverRatio)

	// how many of those new judges should be not judged according to how many
	// poems they've had published and the quality of those poems? (i.e. how
	// many "underdogs")
	numUnderdogs := int(float64(numNewJudges) * (1 - a.Guidelines.Pretension))

	// and finally, how many should be "high brow, zietgiesty" poets
	numPedigreedPoets := numNewJudges - numUnderdogs

	// kick off old judges randomly from committee
	judges := []*types.Poet{}
	keepIdxs, err := utils.NRandomUniqueInts(
		numNewJudges,
		0,
		a.Guidelines.CommitteeSize-1,
	)
	if err != nil {
		a.Error(err.Error())
	}

	// pick the judges to keep
	for _, keepIdx := range keepIdxs {
		judges = append(judges, a.UpcomingIssue.Committee[keepIdx])
	}

	// get underdogs poets
	underdogs, err := types.GetUnderdogPoets(numUnderdogs, a.db, judges...)
	if err != nil {
		a.Error(err.Error())
	}

	// get pedigreed poets
	pedigreedPoets, err := types.GetFancyPoets(
		numPedigreedPoets,
		a.db,
		append(judges, underdogs...)...,
	)
	if err != nil {
		a.Error(err.Error())
	}

	// include all returning and new judges
	judges = append(
		judges,
		append(
			underdogs,
			pedigreedPoets...,
		)...,
	)

	// for each new judge, persist them to the committee!
	for _, judge := range judges {
		err = (&types.IssueCommitteeMembership{
			Poet:  judge,
			Issue: a.UpcomingIssue,
		}).Add(a.db)
	}

	// add these judges to the in-memory issue committee
	a.UpcomingIssue.Committee = judges
}

func (s *MagazineAdministrator) AllowPoetsToLearn() {
	// do stuff
}

func (a *MagazineAdministrator) CleanUp() {
	err := os.RemoveAll(a.limboDir)
	if err != nil {
		a.Error(err.Error())
	}
}

func (a *MagazineAdministrator) UpdateTuneables() {
	// TODO (cw|9.14.2018) pickup new magazine parameters from the environment, or
	// that have been programaticaly set somehow (through an admin API endpoint?)

	// TODO (cw|9.15.2018) update parameters according to randomness/metaRandomness
}
