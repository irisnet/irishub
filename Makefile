#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
NetworkType := $(shell if [ -z ${NetworkType} ]; then echo "mainnet"; else echo ${NetworkType}; fi)
CURRENT_DIR = $(shell pwd)
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
DOCKER := $(shell which docker)

# default mainnet EVM_CHAIN_ID
EVM_CHAIN_ID ?= 6688

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
		  -X github.com/cosmos/cosmos-sdk/version.AppName=iris \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/irisnet/irishub/v2/types.EIP155ChainID=$(EVM_CHAIN_ID) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

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

build: check-evm-chain-id go.sum
ifeq ($(OS),Windows_NT)
	@go build $(BUILD_FLAGS) -o build/iris.exe ./cmd/iris
else
	@go build $(BUILD_FLAGS) -o build/iris ./cmd/iris
endif

build-linux: check-evm-chain-id go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-all-binary: check-evm-chain-id go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build $(BUILD_FLAGS) -o build/iris-linux-amd64 ./cmd/iris
	LEDGER_ENABLED=false GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build $(BUILD_FLAGS)-o build/iris-linux-arm64 ./cmd/iris
	LEDGER_ENABLED=false GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build $(BUILD_FLAGS) -o build/iris-windows-amd64.exe ./cmd/iris

build-contract-tests-hooks:
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests.exe ./cmd/contract_tests
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests ./cmd/contract_tests
endif

install: check-evm-chain-id go.sum
	@go install $(BUILD_FLAGS) ./cmd/iris

check-evm-chain-id:
	@echo "note: EVM_CHAIN_ID is $(EVM_CHAIN_ID)"

update-swagger-docs: statik proto-swagger-gen
	$(BINDIR)/statik -src=lite/swagger-ui -dest=lite -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
    	echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs

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
	rm -rf snapcraft-local.yaml build/ tmp-swagger-gen/

distclean: clean
	rm -rf vendor/

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.13.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

########################################
### Testing


test: test-unit
test-all: test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ${PACKAGES_UNITTEST}

test-sim-nondeterminism-fast:
	@echo "Running non-determinism test..."
	@cd ${CURRENT_DIR}/app && go test -mod=readonly -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@cd ${CURRENT_DIR}/app && $(BINDIR)/runsim -Jobs=4 -SimAppPkg=. -ExitOnFail 50 5 TestAppImportExport

test-sim-after-import: runsim
	@echo "Running application simulation-after-import. This may take several minutes..."
	@cd ${CURRENT_DIR}/app && $(BINDIR)/runsim -Jobs=4 -SimAppPkg=. -ExitOnFail 50 5 TestAppSimulationAfterImport

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" -not -path "*.pb.go" | xargs gofmt -d -s
	go mod verify

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" -not -path "*.pb.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/statik/statik.go" -not -path "*.pb.go" | xargs goimports -w -local github.com/irisnet/irishub/v2

benchmark:
	@go test -mod=readonly -bench=. ./...


########################################
### Local validator nodes using docker and docker-compose

testnet-init:
	@if ! [ -f build/nodecluster/node0/iris/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/home irisnet/irishub iris testnet --v 4 --output-dir /home/nodecluster --chain-id irishub-test --keyring-backend test --starting-ip-address 192.168.10.2 ; fi
	@echo "To install jq command, please refer to this page: https://stedolan.github.io/jq/download/"
	@jq '.app_state.auth.accounts+= [{"@type":"/cosmos.auth.v1beta1.BaseAccount","address":"iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx","pub_key":null,"account_number":"0","sequence":"0"}] | .app_state.bank.balances+= [{"address":"iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx","coins":[{"denom":"uiris","amount":"1000000000000"}]}]' build/nodecluster/node0/iris/config/genesis.json > build/genesis_temp.json ;
	@sudo cp build/genesis_temp.json build/nodecluster/node0/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node1/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node2/iris/config/genesis.json
	@sudo cp build/genesis_temp.json build/nodecluster/node3/iris/config/genesis.json
	@rm build/genesis_temp.json
	@echo "Faucet address: iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx" ;
	@echo "Faucet coin amount: 1000000000000uiris"
	@echo "Faucet key seed: tube lonely pause spring gym veteran know want grid tired taxi such same mesh charge orient bracket ozone concert once good quick dry boss"

testnet-start:
	docker-compose up -d

testnet-stop:
	docker-compose down

testnet-clean:
	docker-compose down
	sudo rm -rf build/*


########################################
### Test ibc nft-transfer
init-golang-rly: kill-dev install
	@echo "Initializing both blockchains..."
	./network/init.sh
	./network/start.sh
	@echo "Initializing relayer..."
	./network/relayer/interchain-nft-config/rly.sh

start: 
	@echo "Starting up test network"
	./network/start.sh

start-rly:
	./network/hermes/start.sh

kill-dev:
	@echo "Killing nftd and removing previous data"
	-@rm -rf ./data
	-@killall nftd 2>/dev/null