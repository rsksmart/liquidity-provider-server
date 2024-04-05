.PHONY: test

tools: download
	go install github.com/parvez3019/go-swagger3@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/conventionalcommit/commitlint@latest
	go env GOPATH
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2
	# installation with brew is because mockery team doesnt recommend to install with go install,
	# if you don't have brew feel free to comment this line and install mockery with other method
	brew install mockery && brew upgrade mockery

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
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=coverage/cover.out.temp ./pkg/... ./internal/... ./cmd/...
	go tool cover -func "coverage/cover.out.temp"
	go tool cover -html="coverage/cover.out.temp"
	rm coverage/cover.out.temp

coverage-report: clean
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=coverage/cover.out ./pkg/... ./internal/... ./cmd/...

test: clean
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=coverage/cover.out.temp  ./pkg/... ./internal/... ./cmd/...
	go tool cover -func "coverage/cover.out.temp"
	rm coverage/cover.out.temp

clean:
	rm -rf build coverage/cover.out.temp
