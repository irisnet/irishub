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
	@go test  -timeout 20m -count 1 -p 1 `go list github.com/irisnet/irishub/client/clitest` -tags=cli_test

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
