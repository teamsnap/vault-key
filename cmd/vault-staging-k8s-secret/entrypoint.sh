#!/bin/sh

set -e

config_auth () {
  if [ -z "${GOOGLE_APPLICATION_CREDENTIALS}" ]; then
    echo "PLUGIN_GOOGLE_APPLICATION_CREDENTIALS environment variable is not set." > /dev/null 2>&1
    exit 1
  else
    location="--zone ${GCLOUD_ZONE}"
    if [ ! -z "${GCLOUD_REGION}" ]
    then
      location="--region ${GCLOUD_REGION}"
    fi

    gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS" && \
    gcloud auth list && \
    gcloud config set account ${FUNCTION_IDENTITY} && \
    gcloud container clusters get-credentials ${CLUSTER_NAME} ${location}
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

config_location () {
  # us-east1
  if [ -z "${GCLOUD_ZONE}" ] && [ -z "${GCLOUD_REGION}" ]; then
    echo "GCLOUD_ZONE and GCLOUD_REGION environment variables are not set. One is required." > /dev/null 2>&1
    exit 1
  elif [ ! -z "${GCLOUD_REGION}" ]
  then
    gcloud config set compute/region ${GCLOUD_REGION}
  else
    gcloud config set compute/zone ${GCLOUD_ZONE}
  fi
}

echo "Initializing gcloud..."
config_project
config_location
config_auth
echo "Finished initializing gcloud!"
./vault-staging-k8s-secret

exit 0
