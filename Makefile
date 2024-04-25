.DEFAULT_GOAL := help

# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))
OUT := cache_clone
PKG := github.com/natemarks/cache_clone
VERSION := 0.0.7
COMMIT := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)


help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

all: run

clean-venv: ## re-create virtual env
	rm -rf .venv
	python3 -m venv .venv
	( \
       source .venv/bin/activate; \
       pip install --upgrade pip setuptools; \
    )

build: ## build the binaries with commit IDs
	env GOOS=linux GOARCH=amd64 \
	go build  -v -o build/$(COMMIT)/${OUT}_linux_amd64 \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${COMMIT}" ${PKG}
	env GOOS=darwin GOARCH=amd64 \
	go build  -v -o build/$(COMMIT)/${OUT}_darwin_amd64 \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${COMMIT}" ${PKG}

release:  ## Build release versions
	mkdir -p build/$(VERSION)/linux/amd64
	env GOOS=linux GOARCH=amd64 \
	go build  -v -o build/$(VERSION)/linux/amd64/${OUT} \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${VERSION}" ${PKG}
	mkdir -p build/$(VERSION)/darwin/amd64
	env GOOS=darwin GOARCH=amd64 \
	go build  -v -o build/$(VERSION)/darwin/amd64/${OUT} \
	-ldflags="-X github.com/natemarks/cache_clone/version.Version=${VERSION}" ${PKG}

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

static: vet lint test

run: server
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*


bump: clean-venv  ## bump version in main branch
ifeq ($(CURRENT_BRANCH), $(MAIN_BRANCH))
	( \
	   source .venv/bin/activate; \
	   pip install bump2version; \
	   bump2version $(part); \
	)
else
	@echo "UNABLE TO BUMP - not on Main branch"
	$(info Current Branch: $(CURRENT_BRANCH), main: $(MAIN_BRANCH))
endif


.PHONY: run build release static upload vet lint