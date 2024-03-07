.PHONY: test

tools: download
	go install github.com/parvez3019/go-swagger3@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/conventionalcommit/commitlint@latest
	go env GOPATH
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

download:
	go mod download

lint:
	test -z $(shell gofmt -l .)
	golangci-lint run -v ./...

validation: lint
	go mod verify
	govulncheck ./... # should fail on non informational vulnerabilities

COMMIT_HASH ?= $(shell git rev-parse HEAD)
SOURCE_VERSION := $(COMMIT_HASH)
build: download
	mkdir -p build && cd build
	@echo "Building liquidity-provider-server $(SOURCE_VERSION)"
	CGO_ENABLED=0 go build -v -installsuffix 'static' -ldflags="-s -X 'main.BuildVersion=$(SOURCE_VERSION)' -X 'main.BuildTime=$(shell date)'" -o ./build/liquidity-provider-server ./cmd/application/main.go

api:
	go-swagger3 --module-path . \
	--main-file-path ./cmd/application/main.go \
	--handler-path ./internal/adapters/entrypoints/rest/handlers \
	--output OpenApi.yml \
	--schema-without-pkg \
	--generate-yaml true

coverage: clean
	mkdir coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=coverage/cover.out ./pkg/... ./internal/... ./cmd/...
	go tool cover -func "coverage/cover.out"
	go tool cover -html="coverage/cover.out"
	rm coverage/cover.out

coverage-report: clean
	mkdir coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=coverage/cover.out ./pkg/... ./internal/... ./cmd/...

test: clean
	mkdir coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=coverage/cover.out  ./pkg/... ./internal/... ./cmd/...
	go tool cover -func "coverage/cover.out"
	rm coverage/cover.out

clean:
	rm -rf build coverage
