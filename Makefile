MODULE   = $(shell env GO111MODULE=on $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))
BIN      = $(CURDIR)/.bin

GOLANGCI_VERSION = v1.47.2

GO           = go
TIMEOUT_UNIT = 5m
TIMEOUT_E2E  = 20m
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m🐱\033[0m")

export GO111MODULE=on

COMMANDS=$(patsubst cmd/%,%,$(wildcard cmd/*))
BINARIES=$(addprefix bin/,$(COMMANDS))

.PHONY: all
all: fmt $(BINARIES) | $(BIN) ; $(info $(M) building executable…) @ ## Build program binary

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Format Go source code
	$Q $(GO) fmt ./...

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)…)
	$Q tmp=$$(mktemp -d); cd $$tmp; \
		env GO111MODULE=on GOPATH=$$tmp GOBIN=$(BIN) $(GO) install $(PACKAGE) \
		|| ret=$$?; \
		env GO111MODULE=on GOPATH=$$tmp GOBIN=$(BIN) $(GO) clean -modcache \
        || ret=$$?; \
		cd - ; \
	  	rm -rf $$tmp ; exit $$ret

FORCE:

bin/%: cmd/% FORCE
	$Q $(GO) build -mod=vendor $(LDFLAGS) -v -o $@ ./$<

KO = $(or ${KO_BIN},${KO_BIN},$(BIN)/ko)
$(BIN)/ko: PACKAGE=github.com/google/ko@latest

.PHONY: apply
apply: | $(KO) ; $(info $(M) ko apply -R -f config/) @ ## Apply config to the current cluster
	$Q $(KO) apply -R -f config

.PHONY: resolve
resolve: | $(KO) ; $(info $(M) ko resolve -R -f config/) @ ## Resolve config to the current cluster
	$Q $(KO) resolve --push=false --oci-layout-path=$(BIN)/oci -R -f config


# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf bin
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:

	@echo $(VERSION)

.PHONY: deploy_tekton
deploy_tekton: clean_tekton | ; $(info $(M) deploying tekton on local cluster …) @ ## Deploying tekton on local clustert
	-kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
	-ko apply -f config;

.PHONY:  clean_tekton
clean_tekton: | ; $(info $(M) deleting tekton from local cluster …) @ ## Deleting tekton on local cluster
	-ko delete -f config;