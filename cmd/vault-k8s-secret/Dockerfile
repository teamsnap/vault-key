FROM golang:1.23.2-bookworm

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

WORKDIR /go/src/github.com/teamsnap/vault-key

# Copy local dependency first for layer caching
COPY pkg/k8s/ pkg/k8s/

# Copy the command module
COPY cmd/vault-k8s-secret/ cmd/vault-k8s-secret/

WORKDIR /go/src/github.com/teamsnap/vault-key/cmd/vault-k8s-secret

RUN go mod download

RUN GONOSUMCHECK=* GOOS=linux GOARCH=amd64 go build -mod=mod -o vault-k8s-secret

RUN useradd --create-home --uid 10001 appuser \
 && chown -R appuser:appuser /go/src/github.com/teamsnap/vault-key
USER appuser

CMD ["./entrypoint.sh"]
