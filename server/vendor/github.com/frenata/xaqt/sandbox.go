package xaqt

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	uuid "github.com/satori/go.uuid"
)

const (
	TmpDirPrefix = "xaqt-"
)

// prepares for execution of user code by creating a temp directory for
// code and input.
//
type sandbox struct {
	// sandbox id (uuidV4)
	ID       string
	language ExecutionDetails
	code     Code
	stdin    string
	options  options
	// docker client connection
	docker *client.Client
	// wait channel for successful container exit
	waitChan <-chan container.ContainerWaitOKBody
	// error channel for container error
	errChan <-chan error
}

// describes how the user Code is represented/structured. either it is in memory as a string,
// or it has been persisted to a file. additionally, the user can include resource files.
//
type Code struct {
	IsFile            bool     // is the code in a file or a string
	String            string   // the code represented as a string
	SourceFileName    string   // the name of the src file the code has been written to
	ResourceFileNames []string // file names for resources used by source file
	Path              string   // the path to the src file and possibly resource files
}

// constructs a new sandbox given...
//
func newSandbox(l ExecutionDetails, code Code, stdin string, opts options) (*sandbox, error) {
	var (
		s   *sandbox
		err error
	)

	// set the API version to use in an environment variable
	// TODO it would be nice to configure based on the docker version
	// a user currently has.... not enough time right now so skipping that.
	err = os.Setenv("DOCKER_API_VERSION", "1.35")
	if err != nil {
		return nil, err
	}

	// init a docker api client
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		// this could occur if docker has not been installed or started
		return nil, err
	}

	// TODO (cw|4.29.2018) if we are spinning up this sandbox from within another docker
	// container, we may want to define a bridge network between them (since they will be
	// sibling containers). I don't know if this is entirely necessary though...
	// THIS NETWORK SETUP SHOULD ACTUALLY GO IN A HIGHER SCOPE (within the struct which
	// actually constructs sandboxes). this way we aren't creating and destroying docker
	// networks all over the place. instead we should check to see if one has been created.

	// define unique network name
	// networkName := fmt.Sprintf("xaqt.%s", uuid.NewV4().String())

	// setup container bridge network if one doesn't already exist.
	// _, err = dockerClient.NetworkCreate(
	// 	context.TODO(),
	// 	networkName,
	// 	types.NetworkCreate{},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	s = &sandbox{
		ID:       uuid.NewV4().String(),
		language: l,
		code:     code,
		stdin:    stdin,
		options:  opts,
		docker:   dockerClient,
	}

	return s, nil
}

// runs user code within the sandbox after preparing the execution environment.
//
func (s *sandbox) run() (string, error) {
	var (
		output string
		err    error
	)

	err = s.prepare()
	if err != nil {
		return "", err
	}

	output, err = s.execute()
	if err != nil {
		return "", err
	}

	return output, nil
}

// prepares the execution environment and the sandbox docker container.
//
func (s *sandbox) prepare() error {
	var err error

	err = s.PrepareTmpDir()
	if err != nil {
		return err
	}

	err = s.PrepareContainer()
	if err != nil {
		return err
	}

	return nil
}

// prepares the execution environment by copying all resources (user code, input files,
// and execution payload) into a temporary directory.
//
func (s *sandbox) PrepareTmpDir() error {
	// create tmp directory for keeping all code, inputs, and results
	tmpFolder, err := ioutil.TempDir(s.options.execDir, TmpDirPrefix)
	if err != nil {
		return err
	}

	// modify perms on tmp dir
	if err := os.Chmod(tmpFolder, 0777); err != nil {
		return err
	}

	// record tmpdir for easy deletion
	s.options.execDir = tmpFolder

	// ensure that execMountDir has been set
	if s.options.execMountDir == "" {
		s.options.execMountDir = filepath.Dir(s.options.execDir)
	}

	// write source file and possibly resource files into tmp dir
	switch s.code.IsFile {
	case true:
		// copy source file into tmp dir
		err = s.copyFile(
			filepath.Join(s.code.Path, s.code.SourceFileName),
			filepath.Join(tmpFolder, s.code.SourceFileName),
		)
		if err != nil {
			return err
		}

		// copy resource files into tmp dir
		for _, ResourceFileName := range s.code.ResourceFileNames {
			err = s.copyFile(
				filepath.Join(s.code.Path, ResourceFileName),
				filepath.Join(tmpFolder, ResourceFileName),
			)
			if err != nil {
				return err
			}
		}
	case false:
		// write source file into tmp dir
		err = ioutil.WriteFile(tmpFolder+"/"+s.language.SourceFile, []byte(s.code.String), 0777)
		if err != nil {
			return err
		}
	}

	// write a file for stdin
	err = ioutil.WriteFile(tmpFolder+"/inputFile", []byte(s.stdin), 0777)
	if err != nil {
		return err
	}

	return nil
}

// create docker container for running code and stream container's stdout to our stdout.
//
func (s *sandbox) PrepareContainer() error {
	var (
		ctx         = context.Background()
		srcFilename string
		err         error
	)

	// the call to #PrepareTmpdir creates a tmp directory within the specified execDir,
	// we want to mount this tmp dir into the sandbox as /usercode/, so we must append
	// the new execDir suffix directory onto the execMountDir.
	s.options.execMountDir = filepath.Join(
		s.options.execMountDir,
		filepath.Base(s.options.execDir), // get the suffix dir just created in #PrepareTmpDir
	)

	// get source file name
	if s.code.IsFile {
		srcFilename = s.code.SourceFileName
	} else {
		srcFilename = s.language.SourceFile
	}

	// create docker container for executing user code
	_, err = s.docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: s.options.image,
			Cmd: []string{
				"/entrypoint/script.sh",
				s.language.Compiler,
				srcFilename,
				s.language.OptionalExecutable,
				s.language.CompilerFlags,
				s.options.inputType,
			},
			// run the sandbox container as a specific user.
			User: "mysql", // TODO (cw|4.29.2018) change this to a constant
			// StopTimeout:  s.options.timeout, // TODO this needs to be a *int ...
			AttachStdout: true, // TODO (cw|8.21.2018) do we need this?
			AttachStderr: true, // TODO (cw|8.21.2018) do we need this?
			Tty:          true, // TODO (cw|8.21.2018) do we need this?
		},
		&container.HostConfig{
			// remove container from host once it exits
			AutoRemove: true,
			// specify the mount point(s) for the sandbox
			Binds: []string{s.options.execMountDir + ":/usercode"},
		},
		nil, // no network config currently
		s.ID,
	)
	if err != nil {
		return err
	}

	// setup stdout stream from container
	// TODO (cw|8.21.2018) do we need this?
	hijackedResp, err := s.docker.ContainerAttach(
		ctx,
		s.ID,
		types.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		return err
	}

	// start hijacking stdout/stderr
	// TODO (cw|8.21.2018) do we need this?
	go func() {
		defer hijackedResp.Close()

		io.Copy(os.Stdout, hijackedResp.Reader)
	}()

	// setup channels to wait for container to stop and be removed.
	// NOTE (cw|8.21.2018) we need WaitConditionRemoved since it is apparently
	// not enough to just wait for the process to stop. Waiting for the process
	// to merely stop resulted in race conditions between the stdout writer in the
	// container and this parent process...
	s.waitChan, s.errChan = s.docker.ContainerWait(
		context.Background(),
		s.ID,
		container.WaitConditionRemoved,
	)

	return nil
}

// executes user code within the sandbox docker container.
//
// returns TODO (cw|4.29.2018) ???
//
func (s *sandbox) execute() (string, error) {
	var (
		ctx = context.Background()
		err error
	)

	// defer cleanup
	defer s.cleanup()

	// okay lets start the container...
	err = s.docker.ContainerStart(
		ctx,
		s.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return "", err
	}

	// wait for the container to stop
	select {
	case <-s.waitChan:
		// ok. the docker process has stopped and the container has been removed.

		// get the errors file
		errorBytes, err := ioutil.ReadFile(s.options.execDir + "/errors")
		if err != nil {
			// there was an error reading the errors file, perhaps it is missing?
			return "", err
		}

		// did the process error?
		if len(errorBytes) > 0 {
			// the user code which was run in the container errored.
			err = fmt.Errorf("user code error: %s", errorBytes)

			return "", err
		}

		outputBytes, err := ioutil.ReadFile(s.options.execDir + "/completed")
		if err != nil {
			// there was an error reading the completed file, perhaps it is missing?
			return "", err
		}

		// successfully completed

		return string(outputBytes), nil
	case err = <-s.errChan:
		// the damn container errored
		return "", err
	case <-time.After(s.options.timeout):
		// the damn container process timed out
		log.Printf("%s timed out", s.language.Compiler)
		return "", fmt.Errorf("Timed out")
	}
}

func (s *sandbox) cleanup() {
	// optionally rewrite source and resource files
	if true {
		err := s.rewriteUserFiles()
		if err != nil {
			log.Fatalf("unable to rewrite files to %s", s.code.Path)
		}
	}

	// delete temporary directory once we have finished execution
	os.RemoveAll(s.options.execDir)
}

// overwrites the original source and resource files supplied by the user with those from
// the tmp execution directory after execution has successfully completed.
// you may be wondering:
// (1) Why would we ever want to do this?
//        In situations where the code the user supplies modifies itself or its resource files
//        and we want these changes to persist.
// (2) When would this situation ever occur?
//        If the user code will be run multiple times and must maintain state between calls. A
//        perfect example is a machine learning model which needs to update its parameters
//        between subsequent calls.
//
func (s *sandbox) rewriteUserFiles() error {
	var (
		err error
	)

	// only proceed if the source code is in a file
	if !s.code.IsFile {
		return nil
	}

	// TODO (cw|8.23.2018) add some checks here to ensure that the files aren't
	// too large.

	// copy tmp source file into original dir
	err = s.copyFile(
		filepath.Join(s.options.execDir, s.code.SourceFileName),
		filepath.Join(s.code.Path, s.code.SourceFileName),
	)
	if err != nil {
		return err
	}

	// copy tmp resource files into original dir
	for _, ResourceFileName := range s.code.ResourceFileNames {
		err = s.copyFile(
			filepath.Join(s.options.execDir, ResourceFileName),
			filepath.Join(s.code.Path, ResourceFileName),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// short utility for easily copying files
//
func (s *sandbox) copyFile(src, dest string) error {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

// TODO (cw|4.29.2018) this cleanup should be in Context (which is initialized once in the
// calling code)
// func (s *sandbox) CleanUp() error {
// 	// remove the network
// 	err := s.docker.NetworkRemove(context.TODO(), executer.Network)
// 	if err != nil && !client.IsErrNotFound(err) {
// 		// something is very wrong here
// 		panic(err)
// 	}

// }
