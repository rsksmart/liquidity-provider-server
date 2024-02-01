tools: download
	go install github.com/parvez3019/go-swagger3@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go env GOPATH
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

download:
	go mod download

lint:
	golangci-lint run -v ./...

validation: lint
	go mod verify
	govulncheck ./... # should fail on non informational vulnerabilities

build: download
	mkdir -p build && cd build
	CGO_ENABLED=0 go build -v -installsuffix 'static' -ldflags="-s" -o ./build/liquidity-provider-server ./cmd/application/main.go

api:
	go-swagger3 --module-path . \
	--main-file-path ./cmd/application/main.go \
	--handler-path ./internal/adapters/entrypoints/rest/handlers \
	--output OpenApi.yml \
	--schema-without-pkg \
	--generate-yaml true
