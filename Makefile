.PHONY: all test format lint tidy

GOLANGCI_LINT_VERSION := v2.14.0
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all:
	GO111MODULE=on go build .

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test:
	GO111MODULE=on go test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...
	GO111MODULE=on go tool cover -func=coverage.out

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	GO111MODULE=on $(GOLANGCI_LINT) fmt ./...

lint:
	GO111MODULE=on $(GOLANGCI_LINT) run ./...

tidy:
	GO111MODULE=on go mod tidy
