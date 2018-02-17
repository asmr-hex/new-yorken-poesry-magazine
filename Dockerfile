# multi-stage build

# build the client
FROM node:9.4.0-alpine as client

WORKDIR /usr/app/client/
COPY client/package*.json ./
RUN npm install -qy
COPY client/ ./
RUN npm run build


# build the server
FROM golang:1.9-alpine as server

WORKDIR /go/src/github.com/connorwalsh/new-yorken-poesry-magazine/server
COPY ./server .
# no need to go get since we are vendoring all our deps
RUN go install -v


# copy server binary and javascript bundle into final resting place
FROM alpine

WORKDIR /usr/app/
COPY --from=server /go/bin/server .
COPY --from=client /usr/app/client/build/ ./client/build/

ENV PORT 8080
EXPOSE 8080

CMD ["/usr/app/server"]
