FROM golang:alpine

ENV GO111MODULE auto

WORKDIR /go/src/gitlab.com/greenly/go-rest-api

COPY . .

RUN go get -d -v ./...

RUN go install gitlab.com/greenly/go-rest-api

ENTRYPOINT /go/bin/go-rest-api

EXPOSE 8090
