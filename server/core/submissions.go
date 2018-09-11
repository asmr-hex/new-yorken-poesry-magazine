package core

import (
	"database/sql"
	"os"
	"time"
)

type SubmissionService struct {
	*Logger
	Period time.Duration
	db     *sql.DB
}

func NewSubmissions(db *sql.DB) *SubmissionService {
	return &SubmissionService{
		Logger: NewLogger(os.Stdout),
		Period: time.Hour * 24 * 7, // TODO make configurable
		db:     db,
	}
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

		s.ChooseNewCommitteeMembers()
	}
}

func (s *SubmissionService) OpenCallForSubmissions() {
	// var (
	// 	poets []*types.Poet
	// 	err   error
	// )

	// // get all candidate poets
	// poets, err = types.ReadPoets(s.db)
	// if err != nil {
	// 	// this may happen if the database goes down or something...
	// 	// but it shouldn't normally happen  ¯\_(ツ)_/¯
	// 	s.Error(err.Error())
	// }

	// run each candidate poet

	// temporarily persist generated poems
}

func (s *SubmissionService) SelectWinningPoems() {
	// for each poem-committe-member pair, generate score

	// arrive at consensus
}

func (s *SubmissionService) ChooseNewCommitteeMembers() {
	// choose new committee members
}

func (s *SubmissionService) UpdatePoets() {
	// do stuff
}
