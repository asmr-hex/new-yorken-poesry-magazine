package types

import "time"

type Poem struct {
	Id     string
	Date   time.Time
	Author *Poet
	Issue  *Issue
	Likes  Int     // number of likes from readers
	Score  float64 // score assigned by committee
}
