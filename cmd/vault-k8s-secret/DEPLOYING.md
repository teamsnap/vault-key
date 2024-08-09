# Deploying
This is manual for now

## Prerequisites
* gcloud cli and both staging/production project credentials
* docker desktop

## Staging
Replace `x.x.x` version with the new version.

```sh
docker build --platform=linux/amd64 -f Dockerfile -t us.gcr.io/staging-205121/vault-key/vault-key:x.x.x .
docker push us.gcr.io/staging-205121/vault-key/k8-secret:x.x.x
```
## Production
Replace `x.x.x` version with the new version.

```sh
docker build --platform=linux/amd64 -f Dockerfile -t us.gcr.io/production-195315/vault-key/vault-key:x.x.x .
docker push us.gcr.io/production-195315/vault-key/k8-secret:x.x.x
```
