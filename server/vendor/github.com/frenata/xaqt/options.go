package xaqt

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// An Option is a function that performs some kind of configuration on
// the context.
// No current implementations return errors, but it is included so guards
// can be added as desired.
// Idea taken from https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type option func(*Context) error

func DataPath() string {
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Fatal("Fatal: 'GOPATH' is not set, cannot locate the data path.")
	}

	return filepath.Join(gopath, "src/github.com/frenata/xaqt/data/")
}

// defaultOptions provides some useful defaults if the user provides none.
func defaultOptions(c *Context) error {

	c.path = DataPath()

	if runtime.GOOS == "darwin" {
		c.execDir = "/tmp"
	}

	c.image = DEFAULT_DOCKER_IMAGE

	c.timeout = time.Second * 5
	return nil
}

// Timeout configures how long evaluation should run before it is killed.
func Timeout(t time.Duration) option {
	return func(c *Context) error {
		c.timeout = t
		return nil
	}
}

// Image configures which docker image should be used for evaluation.
func Image(i string) option {
	return func(c *Context) error {
		c.image = i
		return nil
	}
}

// Path configures the folder with the execution script and "Payload" dir.
func Path(p string) option {
	return func(c *Context) error {
		c.path = p
		return nil
	}
}

// ExecDir configures where the user code will be executed and where the results
// will reside upon completion. this function will also set the execMountDir if it
// has not been previously set (making the assumption that execDir == execMountDir).
//
func ExecDir(f string) option {
	return func(c *Context) error {
		c.execDir = f

		if c.execMountDir == "" {
			c.execMountDir = c.execDir
		}

		return nil
	}
}

// ExecMountDir configures where the user code execution / results directory is mounted
// on the host. This only needs to be spcified if the main application is run within
// its own docker container which will be spinning up the xaqt sandbox as a sibling
// container using the host's docker daemon.
//
func ExecMountDir(f string) option {
	return func(c *Context) error {
		c.execMountDir = f
		return nil
	}
}
