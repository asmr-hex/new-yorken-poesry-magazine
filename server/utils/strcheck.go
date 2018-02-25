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
	alphaNumerics = "[\\p{L}a-zA-Z0-9]"
	symbols       = "[\\._\\- ]"
	// underscores        = "_[^\\.|_|-| ])?"
	// periods            = "\\.([^\\.|_|-| ])?"
	// hyphens            = "-([^\\.|_|-| ])?"
	// spaces             = " ([^\\.|_|-| ])?"
	usernameRegexprStr = fmt.Sprintf(
		"(^%s(%s|%s){%d,%d}%s$)|(^%s$)",
		alphaNumerics,
		// underscores,
		// periods,
		// hyphens,
		// spaces,
		symbols,
		alphaNumerics,
		MIN_USERNAME_LENGTH-1,
		MAX_USERNAME_LENGTH-1,
		alphaNumerics,
		alphaNumerics,
	)
	usernameRegexpr = regexp.MustCompile(usernameRegexprStr)
)

func ValidateUsername(username string) error {
	isValid := usernameRegexpr.MatchString(username)
	if !isValid {
		return fmt.Errorf("invalid username %s provided", username)
	}

	return nil
}
