FROM golang:1.13.3-stretch

ENV \
 GOARCH="amd64" \
 PATH="$PATH:/usr/local/go/bin" \
 GO111MODULE=on

WORKDIR /go/src/github.com/teamsnap/vault-key

RUN go version \
 && go get -u golang.org/x/lint/golint

CMD bash
