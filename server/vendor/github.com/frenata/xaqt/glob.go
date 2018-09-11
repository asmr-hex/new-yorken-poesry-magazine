package xaqt

import "strings"

// seperator is used to delineate inputs when running code through the sandbox
// seperator must be maintained in the scripts in /payload/ if the values don't match things break
const seperator = "\n*-BRK-*\n"

// glob and unglob combine and seperate groups of input or output (compiler deals with globs of text seperated by separator but outside compilation a []string is preferred)
func glob(stdins []string) string {
	glob := strings.Join(stdins, seperator) + seperator
	return glob
}

func unglob(glob string) []string {
	stdins := strings.Split(glob, seperator)
	return stdins[:len(stdins)-1]
}
