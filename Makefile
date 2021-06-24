GO    := go
PROMU := $(GOPATH)/bin/promu

PREFIX              ?= $(shell pwd)
BIN_DIR             ?= $(shell pwd)/.build

# all: format lint build test

build: $(PROMU)
	@echo ">> building binaries"
	$(PROMU) build --prefix $(PREFIX)

build-linux: $(PROMU)
	@echo ">> building binaries"
	GOOS=linux GOARCH=amd64 $(PROMU) build --prefix $(PREFIX)

crossbuild: $(PROMU)
	@echo ">> crossbuilding binaries"
	$(PROMU) crossbuild -v

tarball: $(PROMU)
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

# deps
promu: $(PROMU)

$(PROMU):
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) get -u github.com/prometheus/promu

.PHONY: all build crossbuild  tarball promu  $(PROMU)
