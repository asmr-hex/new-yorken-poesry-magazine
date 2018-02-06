package types

import "time"

type Poem struct {
	Id     string
	Date   time.Time
	Author *Poet
	Issue  *Issue
	Likes  int     // number of likes from readers
	Score  float64 // score assigned by committee
}
