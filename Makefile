# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
BASE_DIR := $(shell pwd)
UNAME_S := $(shell uname -s)
APP_NAME := vault-key

export
# ===================== [END] Global Variable Declaration ==================== #

# =========================== [START] Build Targets ========================== #
docker_build:
	@docker build -f $(BASE_DIR)/build/package/Dockerfile -t teamsnap/$(APP_NAME):latest .

docker_build_vault_init:
	@docker build -f $(BASE_DIR)/cmd/vault-init/Dockerfile -t teamsnap/$(APP_NAME)/vault-init:latest $(BASE_DIR)/cmd/vault-init

docker_build_vault_k8s_secret:
	@docker build -f $(BASE_DIR)/cmd/vault-k8s-secret/Dockerfile -t teamsnap/$(APP_NAME)/vault-k8s-secret:latest $(BASE_DIR)/cmd/vault-k8s-secret
# ============================ [END] Build Targets =========================== #

# ============================ [START] Run Targets =========================== #
docker_run:
	@docker run -it --rm \
		-v $(BASE_DIR):/go/src/github.com/teamsnap/$(APP_NAME) \
	 teamsnap/$(APP_NAME):latest bash
# ============================= [END] Run Targets ============================ #

# ======================== [START] Formatting Targets ======================== #
gofmt:
	@go fmt github.com/teamsnap/$(APP_NAME)/cmd/...
	@go fmt github.com/teamsnap/$(APP_NAME)/pkg/...

golint:
	@golint github.com/teamsnap/$(APP_NAME)/cmd/...
	@golint github.com/teamsnap/$(APP_NAME)/pkg/...

lint: gofmt golint
# ========================= [END] Formatting Targets ========================= #

# ============================ [START] Test Targets ========================== #
test:
	@go test -v -cover github.com/teamsnap/$(APP_NAME)
# ============================= [END] Test Targets =========================== #

# ======================= [START] Documentation Scripts ====================== #
godoc:
	@godoc -http=":6060"
# ==============-========= [END] Documentation Scripts =========-============= #
