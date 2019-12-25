#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
NetworkType := $(shell if [ -z ${NetworkType} ]; then echo "mainnet"; else echo ${NetworkType}; fi)

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=iris \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=iris \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=iriscli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/irisnet/irishub/cmd/config.NetworkType=${NetworkType} \
		  -X github.com/cosmos/cosmos-sdk/types.reDnmString=[a-z][a-z0-9:-]{2,15}

denomflags = -X github.com/cosmos/cosmos-sdk/types.reDnmString=[a-z][a-z0-9:-]{2,15}

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# The below include contains the tools target.

all: tools install lint

# The below include contains the tools.
include contrib/devtools/Makefile

build: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/iris.exe ./cmd/iris
	go build $(BUILD_FLAGS) -o build/iriscli.exe ./cmd/iriscli
else
	go build $(BUILD_FLAGS) -o build/iris ./cmd/iris
	go build $(BUILD_FLAGS) -o build/iriscli ./cmd/iriscli
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-contract-tests-hooks:
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests.exe ./cmd/contract_tests
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests ./cmd/contract_tests
endif

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/iris
	go install $(BUILD_FLAGS) ./cmd/iriscli

update-swagger-docs: statik
	$(BINDIR)/statik -src=lite/swagger-ui -dest=lite -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
    	echo "\033[92mSwagger docs are in sync\033[0m";\
    fi

install-tool: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iristool

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/iris -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf snapcraft-local.yaml build/

distclean: clean
	rm -rf vendor/

########################################
### Testing


test: test-unit test-build
test-all: test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -ldflags='$(denomflags)' -mod=readonly -tags='ledger test_ledger_mock' ${PACKAGES_UNITTEST}

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-build: build
	@go test -ldflags='$(denomflags)' -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test -v


lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" | xargs gofmt -d -s
	go mod verify

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" | xargs goimports -w -local github.com/irisnet/irishub

benchmark:
	@go test -mod=readonly -bench=. ./...


########################################
### Local validator nodes using docker and docker-compose

testnet-init:
	@if ! [ -f build/iris ]; then $(MAKE) build_linux ; fi
	@if ! [ -f build/nodecluster/node0/iris/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris testnet --v 4 --output-dir /home/nodecluster --chain-id irishub-test --starting-ip-address 192.168.10.2 ; fi
	@echo "To install jq command, please refer to this page: https://stedolan.github.io/jq/download/"
	@if [ ${NetworkType} = "testnet" ]; then jq '.app_state.accounts+= [{"address": "faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm", "coins": [ "1000000iris" ], "sequence_number": "0", "account_number": "0"}]' build/nodecluster/node0/iris/config/genesis.json > build/genesis_temp.json ; else jq '.app_state.accounts+= [{"address": "iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx", "coins": [ "1000000iris" ], "sequence_number": "0", "account_number": "0"}]' build/nodecluster/node0/iris/config/genesis.json > build/genesis_temp.json ; fi
	@sudo cp build/genesis_temp.json build/nodecluster/node0/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node1/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node2/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node3/iris/config/genesis.json
	@rm build/genesis_temp.json
	@if [ ${NetworkType} = "testnet" ]; then echo "Faucet address: faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm" ; else echo "Faucet address: iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx" ; fi
	@echo "Faucet coin amount: 1000000iris"
	@echo "Faucet key seed: tube lonely pause spring gym veteran know want grid tired taxi such same mesh charge orient bracket ozone concert once good quick dry boss"

testnet-start:
	docker-compose up -d

testnet-stop:
	docker-compose down

testnet-clean:
	docker-compose down
	sudo rm -rf build/*