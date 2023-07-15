GOPATH ?= ${HOME}/go
RACE ?= 0
ENVIRONMENT ?= development
VERSION ?= dev
EXT ?=

BINARY_NAME ?= bas-celik

.PHONY: all
all: clean test build

.PHONY: test
test:
ifeq ($(RACE), 1)
	go test ./... -race -covermode=atomic -coverprofile=coverage.txt -timeout 5m
else
	go test ./... -covermode=atomic -coverprofile=coverage.txt -timeout 1m
endif

.PHONY: build
build:
ifeq ($(ENVIRONMENT),production)
	go build -ldflags="-s -w" -o ./bin/$(BINARY_NAME)$(EXT) main.go
else ifeq ($(ENVIRONMENT),development)
	go build -o ./bin/$(BINARY_NAME)$(EXT) main.go
else
	echo "Target ${ENVIRONMENT} is not supported"
endif

.PHONY: install
install:
	mv bin/$(BINARY_NAME) ${GOPATH}/bin

.PHONY: commit
commit:
	git add .
ifneq ($(shell git status --porcelain),)
	git commit --author "github-actions[bot] <github-actions[bot]@users.noreply.github.com>" --message "${MESSAGE}" --no-verify
	git push
endif

.PHONY: tidy
tidy:
	@rm -f go.sum
	@go mod tidy

.PHONY: clean
clean:
	@rm -rf ./bin

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt:
	@gofumpt -l -w .

gosec:
	@gosec ./...
