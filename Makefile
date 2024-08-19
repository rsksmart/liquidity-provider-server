.PHONY: test all clean

COVER_FILE = coverage/cover.out
TEMPORAL_COVER_FILE =$(shell pwd)/coverage/cover.out.temp

filter_coverage_file = grep -v "internal/adapters/dataproviders/rootstock/bindings" $(1) > coverage/temp.txt && mv coverage/temp.txt $(1)

tools: download
	go install github.com/parvez3019/go-swagger3@fef3d30b0707883c389261bf26297eebd10d7216 #v1.0.3
	go install golang.org/x/vuln/cmd/govulncheck@latest
	pip3 install pre-commit && pre-commit install
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2
	brew install mockery && brew upgrade mockery 	# installation with brew is because mockery team doesnt recommend to install with go install, if you don't have brew feel free to comment this line and install mockery with other method

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
	--output OpenApi.yml --schema-without-pkg --generate-yaml true

coverage: clean
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=$(TEMPORAL_COVER_FILE) ./pkg/... ./internal/... ./cmd/...
	$(call filter_coverage_file, $(TEMPORAL_COVER_FILE))
	go tool cover -func "$(TEMPORAL_COVER_FILE)" && go tool cover -html="$(TEMPORAL_COVER_FILE)"
	rm $(TEMPORAL_COVER_FILE)

coverage-report: clean
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=$(COVER_FILE) ./pkg/... ./internal/... ./cmd/...
	$(call filter_coverage_file, $(COVER_FILE))

test: clean
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverpkg=./pkg/...,./internal/...,./cmd/... -coverprofile=$(TEMPORAL_COVER_FILE)  ./pkg/... ./internal/... ./cmd/...
	$(call filter_coverage_file, $(TEMPORAL_COVER_FILE))
	go tool cover -func $(TEMPORAL_COVER_FILE)
	rm $(TEMPORAL_COVER_FILE)

clean:
	rm -rf build $(TEMPORAL_COVER_FILE)
