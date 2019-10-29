# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
BASE_DIR := $(shell pwd)
UNAME_S := $(shell uname -s)
APP_NAME := ruby-vault

ifeq ($(UNAME_S),Darwin)
	GOOS = darwin
else
	GOOS = linux
endif

export
# ===================== [END] Global Variable Declaration ==================== #

# =========================== [START] Build Targets ========================== #
build:
	@env GOOS=$(GOOS) GOARCH=amd64 go build -buildmode c-shared -o vault-gem/lib/native/vault.$(GOOS).so main.go

# on MacOS
build_all: docker_build clean
	@echo "Building for MacOS..."
	@env mkdir -p vault-gem/lib/native
	@env GOOS=$(GOOS) GOARCH=amd64 go build -buildmode c-shared -o vault-gem/lib/native/vault.$(GOOS).so main.go
	@echo "Building for Linux..."
	@docker run -it --rm -v $(BASE_DIR):/app teamsnap/$(APP_NAME)-builder:latest make build
	@cd vault-gem && gem build vault-key.gemspec

clean:
	rm -rf vault-gem/lib/native/vault.*.so
	rm -rf vault-gem/lib/native/vault.*.h

docker_build:
	docker build -f Dockerfile -t teamsnap/$(APP_NAME)-builder:latest .
# ============================ [END] Build Targets =========================== #

# ============================ [START] Run Targets =========================== #
docker_run:
	docker run -it --rm -v $(BASE_DIR):/app teamsnap/$(APP_NAME)-builder:latest bash
# ============================= [END] Run Targets ============================ #
