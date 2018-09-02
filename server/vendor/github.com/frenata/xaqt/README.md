[![Build Status](https://travis-ci.org/frenata/xaqt.svg?branch=master)](https://travis-ci.org/frenata/xaqt)

## What is it? ##
`xaqt` (*ɛksəkjuti*) is a *Docker* based sandbox to run untrusted code and return the output to your app. Users can submit their code in any of the supported languages. The system will test the code in an isolated environment. This way you do not have to worry about untrusted code possibly damaging your server intentionally or unintentionally.

## How does it work? ##

A client submits their code and a languageID to the API. The API then creates a new *Docker* container and runs the code using the compiler/interpreter of that language. The program runs inside a virtual machine with limited resources and has a time-limit for execution (20s by default). Once the output is ready it is returned as a result of the API request. The *Docker* container is destroyed and all the files are deleted from the server.

No two coders have access to each other’s *Docker* or files.

## Installation Instructions ##

### OS X ###

If running on OS X, ensure that the command `gtimeout` is installed, commonly via `brew install coreutils`.

### Building the Docker ###

 1. Install docker as appropriate for your platform.
 2. Run `docker pull frenata/xaqt-sandbox` in project root.

### Building the Server ###

 1. Install the Go toolchain as appropriate for your platform.
 2. Run `go get github.com/frenata/xaqt/...`

### Running the Server ###

 1. Set the desired port for xaqt via the environment variable `XAQT_PORT`.
 2. From project root, run `xaqt`.

## Usage Instructions ##

Included with the xaqt library is a simple REST api server. Two endpoints are exposed by the running server:

 * GET `/languages/` : This will return a JSON list with the available target languages.
 * POST `/evaluate/` : This evaluates code, encoded in a JSON body of the following form:

```
{
    "language": "python",
    "stdins": ["1","2"],
    "code": "import sys\nprint(sys.stdin.read())"
}
```

   Returned is a JSON object that reports success or failure of evaluation, and for each element of `stdins`, what the code has printed to `stdout` for that element.

## Development ##
### Vendoring ##
we currently use [`govendor`](https://github.com/kardianos/govendor) as our vendoring tool for external dependencies.
