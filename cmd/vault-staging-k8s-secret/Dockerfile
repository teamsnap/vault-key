FROM golang:1.17.5-bullseye as builder

WORKDIR /go/src/github.com/teamsnap/vault-staging-k8s-secret
COPY go.* .
RUN go mod download

COPY *.go .
RUN env GOOS=linux GOARCH=amd64 go build -o $GOPATH/bin/vault-staging-k8s-secret

FROM debian:bullseye-20211201-slim

RUN apt-get update \
  && apt-get install -y \
  curl \
  apt-transport-https \
  ca-certificates \
  gnupg \
  && echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
  && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - \
  && apt-get update \
  && apt-get install -y google-cloud-sdk

COPY --from=builder /go/bin/vault-staging-k8s-secret /vault-staging-k8s-secret

COPY entrypoint.sh entrypoint.sh

CMD ./entrypoint.sh
