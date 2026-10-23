.PHONY: lint lint-fix test build security ci release dry-release

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test $(if $(filter 1,$(shell go env CGO_ENABLED)),-race) ./...

build:
	go build ./...

security:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

ci: lint test security

release:
	goreleaser release --clean

dry-release:
	goreleaser release --snapshot --clean
