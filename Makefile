SHELL := /bin/bash

GO=go
GO_BUILD_OPT=-mod=readonly
GO_LIB_SRCS=$(wildcard pkg/*/*.go)

CMD_DIRS=$(wildcard cmd/*)
CMDS=$(subst cmd,bin,$(CMD_DIRS))

RM=rm

GO_INTERFACE_SRCS=pkg/repositories/machines.go pkg/usecase/machines.go
GO_MOCK_SRCS=$(join $(dir $(GO_INTERFACE_SRCS)),$(addprefix mock/,$(notdir $(GO_INTERFACE_SRCS))))

TC_ETCD_ENDPOINTS ?= http://127.0.0.1:2379

.SECONDEXPANSION:
bin/%: $$(wildcard cmd/$$*/*.go) $(GO_LIB_SRCS) go.mod
	$(GO) build $(GO_BUILD_OPT) -o $@ ./cmd/$*

.PHONY: all
all: $(CMDS)

.PHONY: clean
clean:
	-$(RM) -rf bin/*

.SECONDEXPANSION:
%.go: $$(join $$(dir $$@),$$(addprefix ../,$$(notdir $$@)))
	$(GO) generate ./$<

.PHONY: mock
mock: $(GO_MOCK_SRCS)

.PHONY: test
test: mock
	TC_ETCD_ENDPOINTS=$(TC_ETCD_ENDPOINTS) $(GO) test -v ./...