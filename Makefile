SHELL := bash
PACKAGE_NAME := github.com/monetha/mth-core

M = $(shell printf "\033[32;1m▶▶▶▶▶\033[0m")

PKGS ?= $(shell go list ./...)

export GO111MODULE := on

.PHONY: dependencies
dependencies: ; $(info $(M) retrieving dependencies…)
	@echo "$(M2) Installing dependencies..."
	go mod download
	@echo "$(M2) Installing goimports..."
	go install golang.org/x/tools/cmd/goimports
	@echo "$(M2) Installing golint..."
	go install golang.org/x/lint/golint
	@echo "$(M2) Installing staticcheck..."
	go install honnef.co/go/tools/cmd/staticcheck
	@echo "$(M2) Installing protoc-gen-go..."
	go install github.com/golang/protobuf/protoc-gen-go

.PHONY: generate-proto
generate-proto: ; $(info $(M) generating protobuf .go files…)
	@protofiles=$$(find . -name "*.proto") && [ -z "$$protofiles" ] || for d in $$protofiles; do protoc --go_out=paths=source_relative:. $$d; done

.PHONY: test
test: ; $(info $(M) running tests…)
	go test -timeout 60s -race -v $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) formatting the code…)
	@echo "$(M2) formatting files..."
	@gofiles=$$(go list -f {{.Dir}} $(PKGS) | grep -v 'mock|producer/messages') && [ -z "$$gofiles" ] || for d in $$gofiles; do goimports -l -w $$d/*.go; done

.PHONY: lint
lint: ; $(info $(M) running lint tools…)
	@echo "$(M2) checking formatting..."
	@gofiles=$$(go list -f {{.Dir}} $(PKGS) | egrep -v 'mock|producer/messages') && [ -z "$$gofiles" ] || unformatted=$$(for d in $$gofiles; do goimports -l $$d/*.go; done) && [ -z "$$unformatted" ] || (echo >&2 "Go files must be formatted with goimports. Following files has problem:\n$$unformatted" && false)
	@echo "$(M2) checking vet..."
	@gofiles=$$(go list -f {{.Dir}} $(PKGS) | egrep -v 'mock|producer/messages') && [ -z "$$gofiles" ] || go vet $$gofiles
	@echo "$(M2) checking staticcheck..."
	@gofiles=$$(go list -f {{.Dir}} $(PKGS) | egrep -v 'mock|producer/messages') && [ -z "$$gofiles" ] || staticcheck $$gofiles
	@echo "$(M2) checking lint..."
	@$(foreach dir,$(PKGS),golint $(dir);)

.PHONY: build
build: ; $(info $(M) building packages…)
	@go build $(PKGS)
