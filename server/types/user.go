package types

type User struct {
	Id       string
	Username string
	Password string
	Email    string
	Poets    []*Poet
}
