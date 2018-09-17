package xaqt

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Context creates a execution context for evaluating user code.
type Context struct {
	compilers Compilers
	options
}

// Message represents details on success or failure of execution.
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// options to control how the sandbox is executed.
//
// NOTE: execMountDir is useful when we are using xaqt as an imported package in
// another application *and* that application is dockerized. In other words, if the
// main application container is spinning up the xaqt sandbox as a sibling container
// (achievable by mounting the host docker daemon as the application container's docker
// daemon), then one must specify the execMountDir.
// when the main application is dockerized, execMountDir is the location *on the host*
// machine where the execDir is mounted. Since we are using the host docker daemon, we
// need to be able to mount the xaqt /usercode/ dir in the execMountDir. This means that
// the application docker container and the xaqt sandbox docker container share a mounted
// directory on the host (execMountDir), i.e.
//
// -- within dockerized application --
//   VOLUME $execMountDir:$execDir
//
// -- within the xaqt sandbox container --
//   VOLUME $execMountDir:/usercode/
//
// if the main application is *not* dockerized, then execDir == execMountDir.
//
type options struct {
	execDir      string // path to tmp execution dir w/ user code
	execMountDir string // path to tmp execution dir w/ user code on docker host
	image        string // name of docker image to run
	timeout      time.Duration
	inputType    string
}

// Uses default sandbox options.
func newDefault(compilers Compilers) *Context {
	c := &Context{compilers, options{}}
	_ = defaultOptions(c)
	return c
}

// NewContext creates a context from a map of compilers and
// some user provided options.
func NewContext(compilers Compilers, options ...option) (*Context, error) {
	context := newDefault(compilers)

	for _, option := range options {
		err := option(context)
		if err != nil {
			return nil, err
		}
	}

	return context, nil
}

// Evaluate code in a given language and for a set of 'stdin's.
func (c *Context) Evaluate(language string, code Code, stdins []string) ([]string, Message) {
	stdinGlob := glob(stdins)
	results, msg := c.run(language, code, stdinGlob)

	return unglob(results), msg
}

// input is n test calls seperated by newlines
// input and expected MUST end in newlines
func (c *Context) run(language string, code Code, stdinGlob string) (string, Message) {
	// log.Printf("launching new %s sandbox", language)
	// log.Printf("launching sandbox...\nLanguage: %s\nStdin: %sCode: Hidden\n", language, stdinGlob)

	lang, ok := c.compilers[strings.ToLower(language)]
	if !ok || lang.Disabled {
		return "", Message{"error", "language not supported"}
	}

	if !code.IsFile && code.String == "" {
		return "", Message{"error", "no code submitted"}
	}

	sb, err := newSandbox(lang.ExecutionDetails, code, stdinGlob, c.options)
	if err != nil {
		log.Printf("sandbox initialization error: %v", err)
		return "", Message{"error", fmt.Sprintf("%s", err)}
	}

	// run the new sandbox
	output, err := sb.run()
	if err != nil {
		log.Printf("sandbox run error: %v", err)
		return output, Message{"error", fmt.Sprintf("%s", err)}
	}

	splitOutput := strings.SplitN(output, "*-COMPILEBOX::ENDOFOUTPUT-*", 2)
	timeTaken := splitOutput[1]
	result := splitOutput[0]

	return result, Message{"success", "compilation took " + timeTaken + " seconds"}
}

// Languages returns a list of available language names.
func (c *Context) Languages() map[string]CompositionDetails { return c.compilers.availableLanguages() }
