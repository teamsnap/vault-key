#!/bin/sh

config_auth () {
  if [ -z "${GOOGLE_APPLICATION_CREDENTIALS}" ]; then
    echo "PLUGIN_GOOGLE_APPLICATION_CREDENTIALS environment variable is not set." > /dev/null 2>&1
    exit 1
  else
    gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS" && \
    gcloud auth list && \
    gcloud config set account ${FUNCTION_IDENTITY} && \
    gcloud container clusters get-credentials ${CLUSTER_NAME} --zone ${GCLOUD_ZONE}
  fi
}

config_project () {
  if [ -z "${GCLOUD_PROJECT}" ]; then
    echo "GCLOUD_PROJECT environment variable is not set." > /dev/null 2>&1
    exit 1
  else
    gcloud config set project ${GCLOUD_PROJECT}
  fi
}

config_region () {
  # us-east1
  if [ -z "${GCLOUD_REGION}" ]; then
    echo "GCLOUD_REGION environment variable is not set." > /dev/null 2>&1
    exit 1
  else
    gcloud config set compute/region ${GCLOUD_REGION}
  fi
}

config_zone () {
  # us-east1
  if [ -z "${GCLOUD_ZONE}" ]; then
    echo "GCLOUD_ZONE environment variable is not set." > /dev/null 2>&1
    exit 1
  else
    gcloud config set compute/zone ${GCLOUD_ZONE}
  fi
}

echo "Initializing gcloud..."
config_auth && \
config_project && \
config_zone && \
config_region && \
echo "Finished initializing gcloud!" && \
./vault-k8s-secret

exit 0
