FROM golang:1.9-alpine

# GOPATH = /go in the golang image
# also $GOPATH/bin has been added to path

WORKDIR /go/src/github.com/connorwalsh/new-yorken-poesry-magazine/server

# copy server src to WORKDIR in container
COPY . .

# since we need to install a go binary (fresh, an fs watcher for development)
# we need to install git, go get the fs watcher, and delete git to reduce image space
# also install goose for db migrations
RUN apk add --no-cache git \
&& go get github.com/pilu/fresh \
&& go get github.com/connorwalsh/gotest \
&& apk del git

# compile and install server binary within container
RUN go install -v

# run the fs watcher, fresh, to recompile go files on all changes
CMD ["fresh"]
