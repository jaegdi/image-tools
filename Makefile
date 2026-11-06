.DEFAULT_GOAL := help

GO ?= go
GOFMT ?= gofmt
BIN ?= image-tool
PKG ?= ./...
TEST ?=
TAG ?=
REMOTE ?= github
DOMAIN ?=

.PHONY: help build test test-one fmt fmt-check vet clean tag config

help: ## Show available development tasks
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the CLI binary
	$(GO) build -o $(BIN) .

config: ## Generate runtime cluster config: make config DOMAIN=cluster.local
	@test -n "$(DOMAIN)" || (echo "Set DOMAIN, for example: make config DOMAIN=cluster.local"; exit 1)
	build/scripts/generate_config.sh "$(DOMAIN)"

test: ## Run all Go tests
	$(GO) test $(PKG)

test-one: ## Run one test: make test-one TEST=TestName [PKG=./ocrequest]
	@test -n "$(TEST)" || (echo "Set TEST, for example: make test-one TEST=TestName"; exit 1)
	$(GO) test $(PKG) -run '^$(TEST)$$' -count=1

fmt: ## Format Go source files
	$(GOFMT) -w $$(find . -type f -name '*.go' -not -path './vendor/*')

fmt-check: ## Check Go source formatting without changing files
	@test -z "$$($(GOFMT) -l $$(find . -type f -name '*.go' -not -path './vendor/*'))"

vet: ## Run go vet on Go packages
	$(GO) vet $(PKG)

clean: ## Remove the CLI binary built by the build target
	rm -f $(BIN)

tag: ## Tag HEAD and push it: make tag TAG=1.2.3 [REMOTE=github]
	@test -n "$(TAG)" || (echo "Set TAG, for example: make tag TAG=1.2.3"; exit 1)
	@git remote get-url "$(REMOTE)" >/dev/null || (echo "Configure the $(REMOTE) remote or set REMOTE to your GitHub remote"; exit 1)
	git tag -a "$(TAG)" -m "$(TAG)"
	@git push "$(REMOTE)" "refs/tags/$(TAG):refs/tags/$(TAG)" || { \
		echo "Push failed; local tag $(TAG) remains. Retry with: git push $(REMOTE) refs/tags/$(TAG):refs/tags/$(TAG)"; \
		exit 1; \
	}
