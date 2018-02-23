package core

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Sessions struct {
	sync.Mutex
	TokenDuration time.Duration
	TokenToUser   map[string]string
	UserToToken   map[string]string
	TokenLastSeen map[string]time.Time
}

func NewSessions(tokenDuration time.Duration) *Sessions {
	session := &Sessions{
		TokenDuration: tokenDuration,
		TokenToUser:   map[string]string{},
		UserToToken:   map[string]string{},
		TokenLastSeen: map[string]time.Time{},
	}

	// being background token sweeper
	go session.SweepUpTokens()

	return session
}

func (s *Sessions) GetTokenByUser(userID string) string {
	var (
		token  string
		exists bool
	)

	// lock writer (since this will be called concurrently)
	s.Lock()
	defer s.Unlock()

	token, exists = s.UserToToken[userID]
	if !exists {
		// create new session token for user
		token = uuid.NewV4().String()
		s.TokenToUser[token] = userID
		s.UserToToken[userID] = token
	}

	// update the token last seen
	s.TokenLastSeen[token] = time.Now()

	return token
}

func (s *Sessions) GetUserByToken(token string) (string, bool) {
	var (
		userId string
		exists bool
	)

	s.Lock()
	defer s.Unlock()

	userId, exists = s.TokenToUser[token]
	if exists {
		// update the token last seen
		s.TokenLastSeen[token] = time.Now()
	}

	return userId, exists
}

func (s *Sessions) SweepUpTokens() {
	ticker := time.NewTicker(s.TokenDuration)
	for {
		<-ticker.C

		s.ExpireSessions()
	}
}

func (s *Sessions) ExpireSessions() {
	expirationThreshold := time.Now().Add(-1 * s.TokenDuration)

	s.Lock()
	defer s.Unlock()

	// go through each token and check expiration time
	for token, lastSeen := range s.TokenLastSeen {
		if lastSeen.After(expirationThreshold) {
			// the token should not be expired!
			continue
		}

		// this toekn should be expired.

		// remove entry from UserToToken
		delete(s.UserToToken, s.TokenToUser[token])

		// remove entry from TokenToUser
		delete(s.TokenToUser, token)

		// remove entry from TokenLastSeen
		delete(s.TokenLastSeen, token)
	}
}
