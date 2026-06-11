# Stage 1: Build the Go binary
FROM golang:1.23.2-bookworm AS builder

WORKDIR /go/src/github.com/teamsnap/vault-key

COPY pkg/k8s/ pkg/k8s/
COPY pkg/vault/ pkg/vault/
COPY cmd/vault-k8s-secret/ cmd/vault-k8s-secret/

WORKDIR /go/src/github.com/teamsnap/vault-key/cmd/vault-k8s-secret

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o /vault-k8s-secret

# Stage 2: Install only the Google Cloud SDK components needed at runtime
FROM debian:bookworm-slim AS sdk

ENV GCLOUD_VERSION="455.0.0"

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl python3 \
 && rm -rf /var/lib/apt/lists/*

RUN curl -sSL "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz" \
    | tar -xz -C /opt \
 && /opt/google-cloud-sdk/install.sh --quiet --usage-reporting=false --path-update=false \
 && /opt/google-cloud-sdk/bin/gcloud components install gke-gcloud-auth-plugin kubectl --quiet \
 && rm -rf /opt/google-cloud-sdk/.install/.backup \
 && rm -rf /opt/google-cloud-sdk/platform/bundledpythonunix \
 && find /opt/google-cloud-sdk -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null; true \
 && find /opt/google-cloud-sdk -name "*.pyc" -delete 2>/dev/null; true

# Stage 3: Minimal runtime image
FROM debian:bookworm-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates python3 \
 && rm -rf /var/lib/apt/lists/*

COPY --from=sdk /opt/google-cloud-sdk /opt/google-cloud-sdk
COPY --from=builder /vault-k8s-secret /app/vault-k8s-secret
COPY cmd/vault-k8s-secret/entrypoint.sh /app/entrypoint.sh

ENV PATH="/opt/google-cloud-sdk/bin:${PATH}"

RUN useradd --create-home --uid 10001 appuser \
 && chown -R appuser:appuser /app
USER appuser

WORKDIR /app

CMD ["./entrypoint.sh"]
