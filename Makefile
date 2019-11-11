PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation' | grep -v '/tests')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')

export GO111MODULE = on

all: get_tools install build_cur

COMMIT_HASH := $(shell git rev-parse --short HEAD)

InvariantLevel := $(shell if [ -z ${InvariantLevel} ]; then echo "error"; else echo ${InvariantLevel}; fi)
NetworkType := $(shell if [ -z ${NetworkType} ]; then echo "mainnet"; else echo ${NetworkType}; fi)

BUILD_TAGS = netgo

BUILD_FLAGS = -tags "$(BUILD_TAGS)" -ldflags "\
-X github.com/irisnet/irishub/version.GitCommit=${COMMIT_HASH} \
-X github.com/irisnet/irishub/types.InvariantLevel=${InvariantLevel} \
-X github.com/irisnet/irishub/types.NetworkType=${NetworkType}"
LEDGER_ENABLED ?= true

########################################
### Build/Install

ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      BUILD_TAGS += ledger
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
        BUILD_TAGS += ledger
      endif
    endif
  endif
endif

########################################
### Tools & dependencies

echo_bech32_prefix:
	@echo "If you want to build binaries for testnet, please execute \"source scripts/setTestEnv.sh\" first"
	@echo InvariantLevel=${InvariantLevel}
	@echo NetworkType=${NetworkType}

check_tools:
	cd scripts && $(MAKE) check_tools

check_dev_tools:
	cd scripts && $(MAKE) check_dev_tools

update_tools:
	cd scripts && $(MAKE) update_tools

update_dev_tools:
	cd scripts && $(MAKE) update_dev_tools

get_tools:
	cd scripts && $(MAKE) get_tools

get_dev_tools:
	cd scripts && $(MAKE) get_dev_tools

go_mod_cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: tools
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw_deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/irisnet/irishub/cmd/iris -d 2 | dot -Tpng -o dependency-graph.png

########################################
### Generate swagger docs for irislcd
update_irislcd_swagger_docs:
	@statik -src=lite/swagger-ui -dest=lite -f

########################################
### Compile and Install
install: update_irislcd_swagger_docs echo_bech32_prefix go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iris
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iriscli
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/irislcd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/iristool

build_linux: update_irislcd_swagger_docs echo_bech32_prefix go.sum
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/iris ./cmd/iris && \
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/iriscli ./cmd/iriscli && \
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/irislcd ./cmd/irislcd && \
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/iristool ./cmd/iristool

build_cur: update_irislcd_swagger_docs echo_bech32_prefix go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/iris ./cmd/iris  && \
	go build -mod=readonly $(BUILD_FLAGS) -o build/iriscli ./cmd/iriscli && \
	go build -mod=readonly $(BUILD_FLAGS) -o build/irislcd ./cmd/irislcd && \
	go build -mod=readonly $(BUILD_FLAGS) -o build/iristool ./cmd/iristool


########################################
### Testing

test: test_unit test_cli test_lcd test_sim

test_sim: test_sim_modules test_sim_benchmark test_sim_iris_nondeterminism test_sim_iris_fast

test_unit:
	@go test -mod=readonly $(PACKAGES_NOSIMULATION)

test_cli:
	@go test -mod=readonly -timeout 20m -count 1 -p 4 tests/cli/*

test_lcd:
	@go test -mod=readonly `go list github.com/irisnet/irishub/client/lcd`

test_sim_modules:
	@echo "Running individual module simulations..."
	@go test -mod=readonly $(PACKAGES_SIMTEST)

test_sim_benchmark:
	@echo "Running benchmark test..."
	@go test -mod=readonly ./app -run=none -bench=BenchmarkFullIrisSimulation -v -SimulationCommit=true -SimulationNumBlocks=100 -timeout 24h

test_sim_iris_nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly ./app -run TestAppStateDeterminism -v -SimulationEnabled=true -timeout 10m

test_sim_iris_fast:
	@echo "Running quick Iris simulation. This may take several minutes..."
	@go test -mod=readonly ./app -run TestFullIrisSimulation -v -SimulationEnabled=true -SimulationNumBlocks=100 -timeout 24h

test_sim_iris_slow:
	@echo "Running full Iris simulation. This may take awhile!"
	@go test -mod=readonly ./app -run TestFullIrisSimulation -v -SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationVerbose=true -timeout 24h

testnet_init:
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

testnet_start:
	docker-compose up -d

testnet_stop:
	docker-compose down

testnet_clean:
	docker-compose down
	sudo rm -rf build/*

testnet_unsafe_reset:
	docker-compose down
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node0/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node1/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node2/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node3/iris
