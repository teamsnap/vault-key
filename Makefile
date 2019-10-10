# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
BASE_DIR := $(shell pwd)
UNAME_S := $(shell uname -s)
APP_NAME := vault

export
# ===================== [END] Global Variable Declaration ==================== #

# =========================== [START] Build Targets ========================== #
docker_build:
	@docker build -f $(BASE_DIR)/build/package/Dockerfile -t teamsnap/$(APP_NAME):latest .
# ============================ [END] Build Targets =========================== #

# ============================ [START] Run Targets =========================== #
docker_run:
	@docker run -it --rm \
		-v $(BASE_DIR):/go/src/github.com/teamsnap/$(APP_NAME) \
	 teamsnap/$(APP_NAME):latest bash
# ============================= [END] Run Targets ============================ #

# ======================== [START] Formatting Targets ======================== #
gofmt:
	@go fmt github.com/teamsnap/$(APP_NAME)/...

golint:
	@golint github.com/teamsnap/$(APP_NAME)/...

govet:
	@go vet github.com/teamsnap/$(APP_NAME)/...

lint: gofmt golint govet
# ========================= [END] Formatting Targets ========================= #

# ============================ [START] Test Targets ========================== #
test:
	@go test -v -cover github.com/teamsnap/$(APP_NAME)
# ============================= [END] Test Targets =========================== #

# ======================= [START] Documentation Scripts ====================== #
godoc:
	@godoc -http=":6060"
# ==============-========= [END] Documentation Scripts =========-============= #
