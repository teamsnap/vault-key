# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
BASE_DIR := $(shell pwd)
UNAME_S := $(shell uname -s)
APP_NAME := node-vault

ifeq ($(UNAME_S),Darwin)
	GOOS = darwin
else
	GOOS = linux
endif

export
# ===================== [END] Global Variable Declaration ==================== #

# =========================== [START] Build Targets ========================== #
build:
	@env GOOS=$(GOOS) GOARCH=amd64 go build -buildmode c-archive -o vault.$(GOOS).a main.go

# on MacOS
build_all: docker_build clean
	@echo "Building for MacOS..."
	@env GOOS=$(GOOS) GOARCH=amd64 go build -buildmode c-archive -o vault.$(GOOS).a main.go
	@echo "Building for Linux..."
	@docker run -it --rm -v $(BASE_DIR):/app teamsnap/$(APP_NAME)-builder:latest make build

clean:
	rm -rf vault.*.a
	rm -rf vault.*.h
	rm -rf build

docker_build:
	docker build -f Dockerfile -t teamsnap/$(APP_NAME)-builder:latest .

node_gyp:
	node-gyp configure
	node-gyp build
# ============================ [END] Build Targets =========================== #

# ============================ [START] Run Targets =========================== #
docker_run:
	docker run -it --rm -v $(BASE_DIR):/app teamsnap/$(APP_NAME)-builder:latest bash
# ============================= [END] Run Targets ============================ #
