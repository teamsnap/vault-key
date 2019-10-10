FROM golang:1.13.1-stretch

ENV \
 GOARCH="amd64" \
 PATH="$PATH:/go/bin" \
 GO111MODULE=on

WORKDIR /go/src/github.com/teamsnap/vault-k8s-secret

COPY go.mod go.mod

RUN go version \
 && go mod download \
 && go mod verify

COPY . .

RUN env GOOS=linux GOARCH=amd64 go build -o vault-k8s-secret

CMD ./vault-k8s-secret
