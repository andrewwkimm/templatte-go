help:
	cat Makefile

################################################################################

build:
	make tidy
	make fix
	make lint
	go build ./...
	make test

ci: bin/.tools
	go mod tidy -diff
	go fix -diff ./...
	bin/golangci-lint run
	go build ./...
	go test -race -shuffle=on -coverprofile=coverage.out ./...

clean:
	rm -rf bin coverage.out

fix:
	go fix ./...

format: bin/.tools
	bin/golangci-lint fmt

lint: bin/.tools
	bin/golangci-lint run --fix

setup:
	@if [ "$$(go env GOTOOLCHAIN)" = local ]; then go env -w GOTOOLCHAIN=auto; fi
	@if [ -f bootstrap.go ]; then go run bootstrap.go; fi
	make tools
	bin/lefthook install

test:
	go test -race -shuffle=on -coverprofile=coverage.out ./...

tidy:
	go mod tidy

tools: bin/.tools bin/gopls

vuln: bin/.tools
	bin/govulncheck ./...

################################################################################

# Builds the tools pinned in tools/go.mod into bin/, again whenever the pins change.
bin/.tools: tools/go.mod tools/go.sum
	GOBIN=$(CURDIR)/bin go install -modfile=tools/go.mod tool
	touch $@

# gopls has its own module because it pins a golang.org/x/tools version that
# golangci-lint's requirements would otherwise override.
bin/gopls: tools/gopls/go.mod tools/gopls/go.sum
	GOBIN=$(CURDIR)/bin go install -modfile=tools/gopls/go.mod tool
	touch $@

.PHONY: \
	build \
	ci \
	clean \
	fix \
	format \
	help \
	lint \
	setup \
	test \
	tidy \
	tools \
	vuln
