# Stage 1: Build the Go binary
FROM golang:1.23.2-bookworm AS builder

WORKDIR /go/src/github.com/teamsnap/vault-key

COPY pkg/k8s/ pkg/k8s/
COPY pkg/vault/ pkg/vault/
COPY cmd/vault-k8s-secret/ cmd/vault-k8s-secret/

WORKDIR /go/src/github.com/teamsnap/vault-key/cmd/vault-k8s-secret

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -ldflags="-s -w" -o /vault-k8s-secret

# Stage 2: Download standalone kubectl and gke-gcloud-auth-plugin
FROM debian:bookworm-slim AS tools

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl \
 && rm -rf /var/lib/apt/lists/*

RUN KUBECTL_VERSION=$(curl -sSL https://dl.k8s.io/release/stable.txt) \
 && curl -sSLo /usr/local/bin/kubectl "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" \
 && chmod +x /usr/local/bin/kubectl

RUN curl -sSLo /usr/local/bin/gke-gcloud-auth-plugin \
    "https://dl.google.com/dl/cloudsdk/channels/rapid/components/google-cloud-sdk-gke-gcloud-auth-plugin-linux-x86_64.tar.gz" \
 && chmod +x /usr/local/bin/gke-gcloud-auth-plugin || true

# If the above URL doesn't work, extract from the SDK
RUN if [ ! -s /usr/local/bin/gke-gcloud-auth-plugin ] || [ "$(wc -c < /usr/local/bin/gke-gcloud-auth-plugin)" -lt 1000000 ]; then \
      apt-get update && apt-get install -y --no-install-recommends python3 && \
      curl -sSL "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-455.0.0-linux-x86_64.tar.gz" \
        | tar -xz -C /opt && \
      /opt/google-cloud-sdk/bin/gcloud components install gke-gcloud-auth-plugin --quiet && \
      cp /opt/google-cloud-sdk/bin/gke-gcloud-auth-plugin /usr/local/bin/gke-gcloud-auth-plugin && \
      rm -rf /opt/google-cloud-sdk; \
    fi \
 && chmod +x /usr/local/bin/gke-gcloud-auth-plugin

# Stage 3: Minimal runtime image — no gcloud SDK, no python
FROM alpine:3.20

RUN adduser -D -u 10001 appuser

RUN apk add --no-cache ca-certificates curl jq openssl

COPY --chown=appuser:appuser --from=tools /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --chown=appuser:appuser --from=tools /usr/local/bin/gke-gcloud-auth-plugin /usr/local/bin/gke-gcloud-auth-plugin
COPY --chown=appuser:appuser --from=builder /vault-k8s-secret /app/vault-k8s-secret
COPY --chown=appuser:appuser cmd/vault-k8s-secret/entrypoint.sh /app/entrypoint.sh

ENV USE_GKE_GCLOUD_AUTH_PLUGIN=True

USER appuser
WORKDIR /app

CMD ["./entrypoint.sh"]
