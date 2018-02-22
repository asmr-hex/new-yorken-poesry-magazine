package core

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Sessions struct {
	sync.Mutex
	TokenToUser   map[string]string
	UserToToken   map[string]string
	TokenLastSeen map[string]time.Time
}

func NewSessions() *Sessions {
	return &Sessions{
		TokenToUser:   map[string]string{},
		UserToToken:   map[string]string{},
		TokenLastSeen: map[string]time.Time{},
	}
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

// TODO we eventually want to have a go-routine constantly running in the background
// at a specified interval which will expire and evict session tokens if no requests
// have been made by a user in a certain time window.
func (s *Sessions) ExpireSessions() {

}
