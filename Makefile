#!/usr/bin/make -f

SIMAPP = ./simapp
BINDIR ?= $(GOPATH)/bin
CURRENT_DIR = $(shell pwd)

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.0.0-rc8
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

ldflags = -X github.com/cosmos/cosmos-sdk/types.reDnmString=[a-zA-Z][a-zA-Z0-9/:]{2,127}

all: tools lint

# The below include contains the tools.
include contrib/devtools/Makefile

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

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -ldflags '$(ldflags)' ${PACKAGES_UNITTEST}

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@cd ${CURRENT_DIR}/simapp && go test -mod=readonly -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-nondeterminism-fast:
	@echo "Running non-determinism test..."
	@cd ${CURRENT_DIR}/simapp && go test -mod=readonly -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@cd ${CURRENT_DIR}/simapp && $(BINDIR)/runsim -Jobs=4 -SimAppPkg=. -ExitOnFail 50 5 TestAppImportExport

test-sim-after-import: runsim
	@echo "Running application simulation-after-import. This may take several minutes..."
	@cd ${CURRENT_DIR}/simapp && $(BINDIR)/runsim -Jobs=4 -SimAppPkg=. -ExitOnFail 50 5 TestAppSimulationAfterImport

test-sim-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.iris/config/genesis.json will be used."
	@$(BINDIR)/runsim -Genesis=${HOME}/.iris/config/genesis.json -SimAppPkg=$(SIMAPP) -ExitOnFail 400 5 TestFullAppSimulation

test-sim-multi-seed-long: runsim
	@echo "Running long multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 10 TestFullAppSimulation

test-sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Period=1 -Commit=true -Seed=57 -v -timeout 24h

lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs gofmt -d -s
	go mod verify

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" -not -path "*.pulsar.go" | xargs goimports -w -local github.com/irisnet/irismod

benchmark:
	@go test -mod=readonly -bench=. ./...

###############################################################################
###                        Compile Solidity Contracts                       ###
###############################################################################

CONTRACTS_DIR := contracts
COMPILED_DIR := $(CONTRACTS_DIR)/compiled_contracts

# Compile and format solidity contracts for the erc20 module. Also install
# openzeppeling as the contracts are build on top of openzeppelin templates.
contracts-compile: contracts-clean dep-install create-contracts-abi

# Install openzeppelin solidity contracts
dep-install:
	@echo "Importing openzeppelin contracts..."
	@npm install
	 
# Clean tmp files
contracts-clean:
	@rm -rf node_modules

# Compile, filter out and format contracts into the following format.
create-contracts-abi:
	solc --combined-json abi,bin --evm-version paris --include-path node_modules --base-path $(CONTRACTS_DIR)/  $(CONTRACTS_DIR)/Token.sol | jq '.contracts["Token.sol:Token"]' > $(COMPILED_DIR)/Token.json \
    && solc --combined-json abi,bin --evm-version paris --include-path node_modules --base-path $(CONTRACTS_DIR)/  $(CONTRACTS_DIR)/TokenProxy.sol | jq '.contracts["TokenProxy.sol:TokenProxy"]' > $(COMPILED_DIR)/TokenProxy.json \
	&& solc --combined-json abi,bin --evm-version paris --include-path node_modules --base-path $(CONTRACTS_DIR)/  $(CONTRACTS_DIR)/UpgradeableBeacon.sol | jq '.contracts["UpgradeableBeacon.sol:UpgradeableBeacon"]' > $(COMPILED_DIR)/UpgradeableBeacon.json \


