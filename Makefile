PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation' | grep -v '/prometheus' | grep -v '/clitest' | grep -v '/lcd' | grep -v '/protobuf')
PACKAGES_MODULES=$(shell go list ./... | grep 'modules')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')

all: get_tools get_vendor_deps install

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -ldflags "-X github.com/irisnet/irishub/version.GitCommit=${COMMIT_HASH}"

########################################
### Tools & dependencies

check_tools:
	cd deps_tools && $(MAKE) check_tools

check_dev_tools:
	cd deps_tools && $(MAKE) check_dev_tools

update_tools:
	cd deps_tools && $(MAKE) update_tools

update_dev_tools:
	cd deps_tools && $(MAKE) update_dev_tools

get_tools:
	cd deps_tools && $(MAKE) get_tools

get_dev_tools:
	cd deps_tools && $(MAKE) get_dev_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

draw_deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/irisnet/irishub/cmd/iris -d 2 | dot -Tpng -o dependency-graph.png

########################################
### Generate swagger docs for irislcd
update_irislcd_swagger_docs:
	@statik -src=client/lcd/swaggerui -dest=client/lcd -f

########################################
### Compile and Install
install: update_irislcd_swagger_docs
	go install $(BUILD_FLAGS) ./cmd/iris
	go install $(BUILD_FLAGS) ./cmd/iriscli
	go install $(BUILD_FLAGS) ./cmd/irislcd
	go install $(BUILD_FLAGS) ./cmd/irismon

install_debug:
	go install ./cmd/irisdebug

build_linux: update_irislcd_swagger_docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iris ./cmd/iris && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iriscli ./cmd/iriscli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/irislcd ./cmd/irislcd && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/irismon ./cmd/irismon

build_cur: update_irislcd_swagger_docs
	go build -o build/iris ./cmd/iris  && \
	go build -o build/iriscli ./cmd/iriscli && \
	go build -o build/irislcd ./cmd/irislcd && \
	go build -o build/irismon ./cmd/irismon

build_examples: update_irislcd_swagger_docs
	go build  -o build/iris1 ./examples/irishub1/cmd/iris1
	go build  -o build/iriscli1 ./examples/irishub1/cmd/iriscli1
	go build  -o build/iris2-bugfix ./examples/irishub-bugfix-2/cmd/iris2-bugfix
	go build  -o build/iriscli2-bugfix ./examples/irishub-bugfix-2/cmd/iriscli2-bugfix


install_examples: update_irislcd_swagger_docs
	go install ./examples/irishub1/cmd/iris1
	go install ./examples/irishub1/cmd/iriscli1
	go install ./examples/irishub-bugfix-2/cmd/iris2-bugfix
	go install ./examples/irishub-bugfix-2/cmd/iriscli2-bugfix


build_example_linux: update_irislcd_swagger_docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris1 ./examples/irishub1/cmd/iris1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli1 ./examples/irishub1/cmd/iriscli1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris2-bugfix ./examples/irishub-bugfix-2/cmd/iris2-bugfix
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli2-bugfix ./examples/irishub-bugfix-2/cmd/iriscli2-bugfix

########################################
### Testing

test: test_unit test_cli test_lcd test_sim

test_sim: test_sim_modules test_sim_benchmark test_sim_iris_nondeterminism test_sim_iris_fast

test_unit:
	#@go test $(PACKAGES_NOSIMULATION)
	@go test $(PACKAGES_MODULES)

test_cli:
	@go test  -timeout 20m -count 1 -p 1 client/clitest/utils.go client/clitest/bank_test.go client/clitest/distribution_test.go client/clitest/gov_test.go client/clitest/iparam_test.go client/clitest/irismon_test.go client/clitest/record_test.go client/clitest/service_test.go client/clitest/stake_test.go

test_upgrade_cli:
	@go test  -timeout 20m -count 1 -p 1 client/clitest/utils.go client/clitest/bank_test.go

test_lcd:
	@go test `go list github.com/irisnet/irishub/client/lcd`

test_sim_modules:
	@echo "Running individual module simulations..."
	@go test $(PACKAGES_SIMTEST)

test_sim_benchmark:
	@echo "Running benchmark test..."
	@go test ./app -run=none -bench=BenchmarkFullIrisSimulation -v -SimulationCommit=true -SimulationNumBlocks=100 -timeout 24h

test_sim_iris_nondeterminism:
	@echo "Running nondeterminism test..."
	@go test ./app -run TestAppStateDeterminism -v -SimulationEnabled=true -timeout 10m

test_sim_iris_fast:
	@echo "Running quick Iris simulation. This may take several minutes..."
	@go test ./app -run TestFullIrisSimulation -v -SimulationEnabled=true -SimulationNumBlocks=100 -timeout 24h

test_sim_iris_slow:
	@echo "Running full Iris simulation. This may take awhile!"
	@go test ./app -run TestFullIrisSimulation -v -SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationVerbose=true -timeout 24h

testnet_init:
	@if ! [ -f build/iris ]; then $(MAKE) build_linux ; fi
	@if ! [ -f build/nodecluster/node0/iris/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris testnet --v 4 --output-dir /home/nodecluster --chain-id irishub-test --starting-ip-address 192.168.10.2 ; fi
	@echo "To install jq command, please refer to this page: https://stedolan.github.io/jq/download/"
	@jq '.app_state.accounts+= [{"address": "faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm", "coins": [{ "denom":"iris-atto","amount": "1000000000000000000000000"}], "sequence_number": "0", "account_number": "0"}]' build/nodecluster/node0/iris/config/genesis.json > build/genesis_temp.json
	@jq '.app_state.stake.pool.loose_tokens="1000600000000000000000000.0000000000"' build/genesis_temp.json > build/genesis_temp1.json
	@sudo cp build/genesis_temp1.json build/nodecluster/node0/iris/config/genesis.json
	@sudo cp build/genesis_temp1.json build/nodecluster/node1/iris/config/genesis.json
	@sudo cp build/genesis_temp1.json build/nodecluster/node2/iris/config/genesis.json
	@sudo cp build/genesis_temp1.json build/nodecluster/node3/iris/config/genesis.json
	@rm build/genesis_temp.json build/genesis_temp1.json
	@echo "Faucet address: faa1ljemm0yznz58qxxs8xyak7fashcfxf5lssn6jm"
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
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node0/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node1/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node2/iris
	@docker run --rm -v $(CURDIR)/build:/home ubuntu:16.04 /home/iris unsafe-reset-all --home=/home/nodecluster/node3/iris