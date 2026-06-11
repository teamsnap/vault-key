#!/bin/sh
set -e

if [ -z "${GOOGLE_APPLICATION_CREDENTIALS}" ]; then
  echo "ERROR: GOOGLE_APPLICATION_CREDENTIALS is not set." >&2
  exit 1
fi

if [ -z "${GCLOUD_PROJECT}" ]; then
  echo "ERROR: GCLOUD_PROJECT is not set." >&2
  exit 1
fi

if [ -z "${CLUSTER_NAME}" ]; then
  echo "ERROR: CLUSTER_NAME is not set." >&2
  exit 1
fi

if [ -z "${GCLOUD_ZONE}" ] && [ -z "${GCLOUD_REGION}" ]; then
  echo "ERROR: GCLOUD_ZONE or GCLOUD_REGION must be set." >&2
  exit 1
fi

if [ -n "${GCLOUD_REGION}" ]; then
  LOCATION="${GCLOUD_REGION}"
  LOCATION_PATH="locations/${LOCATION}"
else
  LOCATION="${GCLOUD_ZONE}"
  LOCATION_PATH="zones/${LOCATION}"
fi

echo "Authenticating service account..."
SA_EMAIL=$(jq -r '.client_email' "${GOOGLE_APPLICATION_CREDENTIALS}")
SA_KEY=$(jq -r '.private_key' "${GOOGLE_APPLICATION_CREDENTIALS}")

JWT_HEADER=$(printf '{"alg":"RS256","typ":"JWT"}' | openssl base64 -e -A | tr '+/' '-_' | tr -d '=')
NOW=$(date +%s)
EXP=$((NOW + 3600))
JWT_PAYLOAD=$(printf '{"iss":"%s","scope":"https://www.googleapis.com/auth/cloud-platform","aud":"https://oauth2.googleapis.com/token","iat":%d,"exp":%d}' \
  "${SA_EMAIL}" "${NOW}" "${EXP}" | openssl base64 -e -A | tr '+/' '-_' | tr -d '=')

SIGNING_INPUT="${JWT_HEADER}.${JWT_PAYLOAD}"
SIGNATURE=$(printf '%s' "${SIGNING_INPUT}" | openssl dgst -sha256 -sign <(printf '%s' "${SA_KEY}") | openssl base64 -e -A | tr '+/' '-_' | tr -d '=')

ACCESS_TOKEN=$(curl -s -X POST https://oauth2.googleapis.com/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&assertion=${SIGNING_INPUT}.${SIGNATURE}" \
  | jq -r '.access_token')

if [ -z "${ACCESS_TOKEN}" ] || [ "${ACCESS_TOKEN}" = "null" ]; then
  echo "ERROR: Failed to obtain access token." >&2
  exit 1
fi

echo "Fetching cluster credentials for ${CLUSTER_NAME}..."
CLUSTER_INFO=$(curl -s -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  "https://container.googleapis.com/v1/projects/${GCLOUD_PROJECT}/${LOCATION_PATH}/clusters/${CLUSTER_NAME}")

ENDPOINT=$(echo "${CLUSTER_INFO}" | jq -r '.endpoint')
CA_CERT=$(echo "${CLUSTER_INFO}" | jq -r '.masterAuth.clusterCaCertificate')

if [ -z "${ENDPOINT}" ] || [ "${ENDPOINT}" = "null" ]; then
  echo "ERROR: Failed to fetch cluster endpoint." >&2
  echo "${CLUSTER_INFO}" >&2
  exit 1
fi

echo "Configuring kubeconfig..."
KUBECONFIG_PATH="${HOME}/.kube/config"
mkdir -p "${HOME}/.kube"

cat > "${KUBECONFIG_PATH}" <<KUBECONFIG
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: ${CA_CERT}
    server: https://${ENDPOINT}
  name: gke_${GCLOUD_PROJECT}_${LOCATION}_${CLUSTER_NAME}
contexts:
- context:
    cluster: gke_${GCLOUD_PROJECT}_${LOCATION}_${CLUSTER_NAME}
    user: gke_${GCLOUD_PROJECT}_${LOCATION}_${CLUSTER_NAME}
  name: gke_${GCLOUD_PROJECT}_${LOCATION}_${CLUSTER_NAME}
current-context: gke_${GCLOUD_PROJECT}_${LOCATION}_${CLUSTER_NAME}
users:
- name: gke_${GCLOUD_PROJECT}_${LOCATION}_${CLUSTER_NAME}
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
      installHint: Install gke-gcloud-auth-plugin
      provideClusterInfo: true
      args: []
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: "${GOOGLE_APPLICATION_CREDENTIALS}"
KUBECONFIG

echo "Authenticated to cluster ${CLUSTER_NAME} at ${ENDPOINT}"
./vault-k8s-secret

exit 0
