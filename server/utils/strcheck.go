package utils

import (
	"fmt"
	"regexp"
)

const (
	MIN_USERNAME_LENGTH = 1
	MAX_USERNAME_LENGTH = 18
)

var (
	alphaNumerics      = "[\\p{L}a-zA-Z0-9]"
	symbols            = "[\\._\\- ]"
	usernameRegexprStr = fmt.Sprintf(
		"(^%s(%s|%s){%d,%d}%s$)|(^%s$)",
		alphaNumerics,
		symbols,
		alphaNumerics,
		MIN_USERNAME_LENGTH-1,
		MAX_USERNAME_LENGTH-1,
		alphaNumerics,
		alphaNumerics,
	)

	// usernames can have spaces, hyphens, underscores, periods
	// and unicode characters (including chinese, japanese,
	// cyrillic, etc.).
	usernameRegexpr = regexp.MustCompile(usernameRegexprStr)

	// very basic check for email addresses. As long as it has
	// an @ and a . before the tld, it is considered valid. To
	// *really* check if it is valid, we just send the user an
	// actual email on registration.
	emailRegexpr = regexp.MustCompile("^.+@.+\\..+$")
)

func ValidateUsername(username string) error {
	isValid := usernameRegexpr.MatchString(username)
	if !isValid {
		return fmt.Errorf("invalid username %s provided", username)
	}

	return nil
}

func ValidateEmail(email string) error {
	isValid := emailRegexpr.MatchString(email)
	if !isValid {
		return fmt.Errorf("invalid email address %s provided", email)
	}

	return nil
}
