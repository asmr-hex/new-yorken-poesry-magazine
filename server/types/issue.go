package types

import "time"

type Issue struct {
	Id           string
	Date         time.Time
	Committee    *ReviewCommittee
	Contributors []*Poet
	Poems        []*Poem
	Description  string
}
