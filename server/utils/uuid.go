package utils

import "regexp"

var (
	uuidRegexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
)

// check to see if a string is a UUID V4
func IsValidUUIDV4(uuid string) bool {
	return uuidRegexp.MatchString(uuid)
}

func ScanForSQLInjection(s string) error {
	// maybe we need this in the future, but we don't need this when using prepared sql statements
	return nil
}
