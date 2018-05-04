APP_NAME = oscli
OS ?= darwin
REPO_PATH ?= github.com/glynternet/oscli
BINARY_NAME ?= $(APP_NAME)-$(OS)
BUILD_FILE ?= ./bin/$(BINARY_NAME)
LATESTVERSION_FILE=latest-version

build:
	docker run --rm -v "$(GOPATH)/src/$(REPO_PATH)":/go/src/$(REPO_PATH) -w /go/src/$(REPO_PATH) golang:latest \
	env GOOS=$(OS) \
	go build -v -o $(BUILD_FILE) ./cmd/$(APP_NAME)

install:
	mv ./bin/* /usr/local/bin/
