PROJECT_NAME := $(shell basename "$(PWD)")
PKG := "github.com/kott/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ./...)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

.PHONY: all dep build clean test coverage coverage-report lint fmt

all: build

lint: ## Lint the files
	@golint -set_exit_status ./...

fmt: ## Format the files
	@go fmt ./...

vet: ## Vet the files
	@go vet ./...

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

coverage: ## Generate global code coverage info
	@gocov test ${PKG_LIST}

coverage-report: ## Generate global code coverage report
	@gocov test ${PKG_LIST} | gocov report

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -o $(PROJECT_NAME) -v ./cmd/api/

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
