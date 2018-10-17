package core

import (
	"sync"
	"time"

	"github.com/connorwalsh/new-yorken-poesry-magazine/server/types"
	uuid "github.com/satori/go.uuid"
)

type Verifier struct {
	sync.Mutex
	TokenDuration       time.Duration
	TokenToUser         map[string]types.User
	EmailAddressToToken map[string]string
	TokenIssued         map[string]time.Time
}

func NewVerifier(tokenDuration time.Duration) *Verifier {
	verifier := &Verifier{
		TokenDuration:       tokenDuration,
		TokenToUser:         map[string]types.User{},
		EmailAddressToToken: map[string]string{},
		TokenIssued:         map[string]time.Time{},
	}

	go verifier.SweepUpTokens()

	return verifier
}

func (v *Verifier) RegisterPendingUser(user types.User) string {
	var (
		token string
	)

	v.Lock()
	defer v.Unlock()

	token = uuid.NewV4().String()

	v.TokenToUser[token] = user
	v.EmailAddressToToken[user.Email] = token

	return token
}

func (v *Verifier) GetTokenByEmailAddress(email string) (string, bool) {
	var (
		token  string
		exists bool
	)

	v.Lock()
	defer v.Unlock()

	token, exists = v.EmailAddressToToken[email]

	return token, exists
}

func (v *Verifier) GetUserByToken(token string) (types.User, bool) {
	var (
		user   types.User
		exists bool
	)

	v.Lock()
	defer v.Unlock()

	user, exists = v.TokenToUser[token]

	return user, exists
}

func (v *Verifier) SweepUpTokens() {
	ticker := time.NewTicker(v.TokenDuration)
	for {
		<-ticker.C

		v.ExpireTokens()
	}
}

func (v *Verifier) ExpireTokens() {
	expirationThreshold := time.Now().Add(-1 * v.TokenDuration)

	v.Lock()
	defer v.Unlock()

	// go through each token and check expiration time
	for token, issuedTime := range v.TokenIssued {
		if issuedTime.After(expirationThreshold) {
			// the token should not be expired!
			continue
		}

		// this toekn should be expired.

		// remove entry from TokenToUser
		delete(v.TokenToUser, token)

		// remove entry from TokenLastSeen
		delete(v.TokenIssued, token)
	}
}
