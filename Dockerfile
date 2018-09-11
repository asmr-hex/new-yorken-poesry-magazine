# multi-stage build

# build the server
FROM golang:1.9-alpine as server

WORKDIR /go/src/github.com/connorwalsh/new-yorken-poesry-magazine/server
COPY ./server .
# no need to go get since we are vendoring all our deps
RUN go install -v


# copy server binary into final resting place
FROM alpine

WORKDIR /usr/app/
COPY --from=server /go/bin/server .

# make a volume where we can store uploaded execs on fs
# TODO (cw|9.11.2018) not sure if this works even with docker-compose
VOLUME /poets

# TODO (cw|9.11.2018) not sure if this works even with docker-compose
ENV PORT 8080
EXPOSE 8080

CMD ["/usr/app/server"]
