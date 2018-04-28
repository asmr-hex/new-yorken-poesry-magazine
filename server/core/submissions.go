package core

import (
	"database/sql"
	"time"
)

type Submissions struct {
	Period time.Duration
	db     *sql.DB
}

func NewSubmissions(db *sql.DB) *Submissions {
	return &Submissions{
		Period: time.Hour * 24 * 7, // TODO make configurable
		db:     db,
	}
}

func (s *Submissions) StartScheduler() {
	ticker := time.NewTicker(s.Period)

	// that's right, we plan on publishing this magazine ~*~*~*F O R E V E R*~*~*~
	// just look at this unqualified loop! It's like staring into the void of perpetual
	// poetical motion, like an unbreakable möbius band, lithe yet oppressive.
	for {
		<-ticker.C

		s.OpenCallForSubmissions()

		s.SelectWinningPoems()

		s.ChooseNewCommitteeMembers()
	}
}

func (s *Submissions) OpenCallForSubmissions() {
	// get all candidate poets

	// run each candidate poet

	// persist generated poems
}

func (s *Submissions) SelectWinningPoems() {
	// for each poem-committe-member pair, generate score

	// arrive at consensus
}

func (s *Submissions) ChooseNewCommitteeMembers() {
	// choose new committee members
}
