FROM golang:1.13.1-stretch

ENV \
 GOARCH="amd64" \
 PATH="$PATH:/usr/local/go/bin" \
 GO111MODULE=on

WORKDIR /go/src/github.com/teamsnap/vault

RUN go version \
 && go get -u golang.org/x/lint/golint

CMD bash
