BINARY               = goinstall
VERSION              = $(shell sed -n 's/var version = "\(v[0-9]*\.[0-9]*\.[0-9]*\)"/\1/p' main.go)
SOURCES              = $(shell find . -name '*.go')
GOPKGS               = $(shell go list ./... | grep -v test)
BUILD_FLAGS          = -v
LDFLAGS              = -w -s

.PHONY: default
default: build.local ## Build binary with local OS and architecture

.PHONY: help
help: ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Makefile for golang projects\n\nUsage:\n  make [target]\n\n\Targets:\n"} /^[a-zA-Z_0-9-]+.*[a-zA-Z_0-9-]*:.*?##/ { printf "  %-22s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) }' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all binaries and temporary files
	rm -rf bin

.PHONY: mod
mod.ci: ## Run go mod tidy with dif to exit with failure instead of fix
	go mod tidy -diff

.PHONY: mod
mod: ## Run go mod tidy
	go mod tidy

.PHONY: fmt
fmt: $(SOURCES) ## Format code
	gofmt -w -s $(SOURCES)
	golines -w -m 120 $(SOURCES)

.PHONY: deps
deps: ## Install required dev dependencies
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint --version
	go install github.com/segmentio/golines@latest
	golines --version

.PHONY: pre-commit
pre-commit: mod fmt lint ## Run all pre-commit targets

.PHONY: pre-commit.install
pre-commit.install: ## Install pre-commit target as git pre-commit hook
	echo "make pre-commit" > .git/hooks/pre-commit
	chmod 755 .git/hooks/pre-commit

.PHONY: test
test: ## Run unit tests
	go test -v -race -cover $(GOPKGS)

.PHONY: lint
lint: ## Run lint with fix
	golangci-lint run -v --fix

.PHONY: lint.ci
lint.ci: ## Run lint with failure exit code instead of fix
	golangci-lint run -v

.PHONY: build.local
build.local: build/$(BINARY) ## Build binary with local OS and architecture

.PHONY: build.linux.amd64
build.linux.amd64: build/linux/amd64/$(BINARY) ## Build binary for linux amd64

.PHONY: build.linux.arm64
build.linux.arm64: build/linux/arm64/$(BINARY) ## Build binary for linux arm64

.PHONY: build.darwin.amd64
build.darwin.amd64: build/darwin/amd64/$(BINARY) ## Build binary for macos amd64

.PHONY: build.darwin.arm64
build.darwin.arm64: build/darwin/arm64/$(BINARY) ## Build binary for macos arm64

build/$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o bin/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

build/linux/amd64/$(BINARY): $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/linux/amd64/$(BINARY) -ldflags "$(LDFLAGS)" .

build/linux/arm64/$(BINARY): $(SOURCES)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/linux/arm64/$(BINARY) -ldflags "$(LDFLAGS)" .

build/darwin/amd64/$(BINARY): $(SOURCES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/darwin/amd64/$(BINARY) -ldflags "$(LDFLAGS)" .

build/darwin/arm64/$(BINARY): $(SOURCES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/darwin/arm64/$(BINARY) -ldflags "$(LDFLAGS)" .

.PHONY: build.package
build.package: build.linux.amd64 build.linux.arm64 build.darwin.amd64 build.darwin.arm64 ## Create all tars with binaries for darwin/linux with arm64/amd64
	tar -C bin/linux/amd64 -czvf $(BINARY)-$(VERSION)-linux-amd64.tar.gz .
	tar -C bin/linux/arm64 -czvf $(BINARY)-$(VERSION)-linux-arm64.tar.gz .
	tar -C bin/darwin/amd64 -czvf $(BINARY)-$(VERSION)-darwin-amd64.tar.gz .
	tar -C bin/darwin/arm64 -czvf $(BINARY)-$(VERSION)-darwin-arm64.tar.gz .
