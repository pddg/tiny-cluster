SHELL := /bin/bash

GO=go
GO_BUILD_OPT=-mod=readonly
GO_LIB_SRCS=$(wildcard pkg/*/*.go)

CMD_DIRS=$(wildcard cmd/*)
CMDS=$(subst cmd,bin,$(CMD_DIRS))

RM=rm

GO_INTERFACE_SRCS=pkg/repositories/machines.go pkg/usecase/machines.go
GO_MOCK_SRCS=$(join $(dir $(GO_INTERFACE_SRCS)),$(addprefix mock/,$(notdir $(GO_INTERFACE_SRCS))))

# Tools managed by gex
PROTOC_GEN_GO=bin/protoc-gen-go
MOCKGEN=bin/mockgen

PROTOC=bin/protoc
PROTOC_VERSION=3.13.0
ifeq "$(OS)" "Windows_NT"
	PROTOC_PKG=protoc-$(PROTOC_VERSION)-win64.zip
else
    UNAME=$(shell uname -s)
ifeq "$(UNAME)" "Linux"
	PROTOC_PKG=protoc-$(PROTOC_VERSION)-linux-x86_64.zip
else
	PROTOC_PKG=protoc-$(PROTOC_VERSION)-osx-x86_64.zip
endif
endif

PROTOC_DOWNLOAD_URL=https://github.com/protocolbuffers/protobuf/releases/download
PROTOC_PKG_URL=$(PROTOC_DOWNLOAD_URL)/v$(PROTOC_VERSION)/$(PROTOC_PKG)

PROTO_SRCS=$(wildcard proto/*.proto)
GO_PB_DIR=./pkg/api/pb
GO_PB_SRCS=$(join $(GO_PB_DIR)/,$(patsubst %.proto,%.pb.go,$(notdir $(PROTO_SRCS))))

TC_ETCD_ENDPOINTS ?= http://127.0.0.1:2379

$(PROTOC_GEN_GO) $(MOCKGEN): tools.go
	go generate ./$<

tmp bin:
	mkdir $@

tmp/$(PROTOC_PKG): tmp
	@echo "--- Downloading protoc..."
	wget -q -O tmp/$(PROTOC_PKG) $(PROTOC_PKG_URL)
	@touch $@

$(PROTOC): tmp/$(PROTOC_PKG) bin
ifeq "$(OS)" "Windows_NT"
	@echo "protoc package could not be unziped from CLI in Windows."
	@echo "Please unarchive it into '.\bin'."
else
	@echo "--- Unarchiving protoc package..."
	unzip -q -o $<
	@rm readme.txt
	@touch $@
	@echo "--- protoc successfully installed"
	$@ --version
endif

.SECONDEXPANSION:
bin/%: $$(wildcard cmd/$$*/*.go) $(GO_LIB_SRCS) go.mod bin
	$(GO) build $(GO_BUILD_OPT) -o $@ ./cmd/$*

.PHONY: all
all: $(CMDS)

.PHONY: clean
clean:
	-$(GO) clean
	-$(RM) -rf bin/*

.SECONDEXPANSION:
%.go: $$(join $$(dir $$@),$$(addprefix ../,$$(notdir $$@))) $(MOCKGEN)
	$(GO) generate ./$<

.PHONY: mock
mock: $(GO_MOCK_SRCS)

.SECONDEXPANSION:
%.pb.go: ./proto/$$(subst .pb.go,.proto,$$(notdir $$@)) $(PROTOC) $(PROTOC_GEN_GO)
	$(PROTOC) \
	--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
	-I=$(dir $<) \
	--go_opt=module=github.com/pddg/tiny-cluster \
	--go_out=. \
	$<

.PHONY: pb
pb: $(GO_PB_SRCS)

.PHONY: test
test: mock pb
	TC_ETCD_ENDPOINTS=$(TC_ETCD_ENDPOINTS) $(GO) test -v ./...