SHELL := /bin/bash

GO=go
GO_BUILD_OPT=
GO_LIB_SRCS=$(wildcard pkg/*/*.go)

CMD_DIRS=$(wildcard cmd/*)
CMDS=$(subst cmd,bin,$(CMD_DIRS))

RM=rm

.SECONDEXPANSION:
bin/%: $$(wildcard cmd/$$*/*.go) $(GO_LIB_SRCS) go.mod
	$(GO) build $(GO_BUILD_OPT) -o $@ ./cmd/$*

.PHONY: all
all: $(CMDS)

.PHONY: clean
clean:
	-$(RM) -rf bin/*
