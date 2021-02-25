#!/usr/bin/make -f

########################################
### Simulations

BINDIR ?= $(GOPATH)/bin
SIMAPP = github.com/irisnet/irishub/app

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.iris/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullIrisSimulation -Genesis=${HOME}/.iris/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running Iris import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 25 5 TestIrisImportExport

test-sim-after-import: runsim
	@echo "Running Iris simulation-after-import. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 25 5 TestIrisSimulationAfterImport

test-sim-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.iris/config/genesis.json will be used."
	@$(BINDIR)/runsim -Jobs=4 -Genesis=${HOME}/.iris/config/genesis.json 400 5 TestFullIrisSimulation

test-sim-multi-seed-long: runsim
	@echo "Running multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 10 TestFullAppSimulation

test-sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true

test-sim-iris-benchmark:
	@echo "Running Iris benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullIrisSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

test-sim-iris-profile:
	@echo "Running Iris benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullIrisSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out

.PHONY: runsim test-sim-iris-nondeterminism test-sim-iris-custom-genesis-fast test-sim-iris-fast sim-iris-import-export \
	test-sim-iris-simulation-after-import test-sim-iris-custom-genesis-multi-seed test-sim-iris-multi-seed \
	test-sim-benchmark-invariants test-sim-iris-benchmark test-sim-iris-profile
