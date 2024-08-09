FROM golang:1.22.6-bookworm

ENV \
 GOARCH="amd64" \
 GCLOUD_VERSION="455.0.0" \
 PATH="$PATH:/go/bin"

RUN apt-get update \
 && apt-get install -y \
         apt-transport-https \
         ca-certificates \
         wget

# install gcloud manually because the latest version ships with a kubectl that does not work for our clusters
RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz \
&& tar -xzf google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz \
&& rm google-cloud-sdk-$GCLOUD_VERSION-linux-x86_64.tar.gz \
&& mv google-cloud-sdk /usr/local \
&& ln -s /usr/local/google-cloud-sdk/bin/gcloud /usr/local/bin/gcloud \
&& gcloud components install gke-gcloud-auth-plugin \
&& gcloud components install kubectl  \
&& ln -s /usr/local/google-cloud-sdk/bin/kubectl.$KUBECTL_VERSION /usr/local/bin/kubectl \
&& apt-get remove -y curl wget

WORKDIR /go/src/github.com/teamsnap/vault-k8s-secret

COPY go.mod go.mod

RUN go version \
 && go mod download \
 && go mod verify

COPY . .

RUN env GOOS=linux GOARCH=amd64 go build -o vault-k8s-secret

CMD ./entrypoint.sh
