M = $(shell printf "\033[32;1m▶▶▶▶▶\033[0m")

TOOLS = golang.org/x/tools/cmd/goimports \
         github.com/Masterminds/glide \
         golang.org/x/lint/golint \
         honnef.co/go/tools/cmd/staticcheck \
         honnef.co/go/tools/cmd/gosimple \
         honnef.co/go/tools/cmd/unused

.PHONY: test
test: vendor fmt-check vet simple unused static-check lint ; $(info $(M) running tests…)
	go test -timeout 20s -race -v $$(glide novendor)

.PHONY: tools
tools: ; $(info $(M) building tools…)
	go get -v $(TOOLS)

.PHONY: vendor
vendor: tools ; $(info $(M) retrieving dependencies…)
	glide install

.PHONY: fmt-check
fmt-check: vendor tools ; $(info $(M) checking formattation…)
	gofiles=$$(go list -f {{.Dir}} $$(glide novendor) | grep -v mock) && [ -z "$$gofiles" ] || unformatted=$$(for d in $$gofiles; do goimports -l $$d/*.go; done) && [ -z "$$unformatted" ] || (echo >&2 "Go files must be formatted with goimports. Following files has problem:\n$$unformatted" && false)

.PHONY: fmt
fmt: vendor tools ; $(info $(M) formatting the code…)
	gofiles=$$(go list -f {{.Dir}} $$(glide novendor) | grep -v mock) && [ -z "$$gofiles" ] || for d in $$gofiles; do goimports -l -w $$d/*.go; done

.PHONY: vet
vet: vendor tools ; $(info $(M) checking correctness of the code…)
	go vet $$(glide novendor)

.PHONY: simple
simple: vendor tools ; $(info $(M) detecting code that could be rewritten in a simpler way…)
	gosimple $$(glide novendor)

.PHONY: unused
unused: vendor tools ; $(info $(M) checking Go code for unused constants, variables, functions and types…)
	unused $$(glide novendor)

.PHONY: static-check
static-check: vendor tools ; $(info $(M) detecting bugs and inefficiencies in code…)
	staticcheck -version $$(glide novendor)

.PHONY: lint
lint: vendor tools ; $(info $(M) running golint…)
	@./ci/run-golint.sh
