FROM golang:1.22.6-bookworm

ENV \
 GOARCH="amd64" \
 PATH="$PATH:/go/bin"

RUN apt-get update \
 && apt-get install -y \
         apt-transport-https \
         ca-certificates \
 && echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
 && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - \
 && apt-get update \
 && apt-get install -y google-cloud-sdk

WORKDIR /go/src/github.com/teamsnap/vault-k8s-secret

COPY go.mod go.mod

RUN go version \
 && go mod download \
 && go mod verify

COPY . .

RUN env GOOS=linux GOARCH=amd64 go build -o vault-k8s-secret

CMD ./entrypoint.sh
