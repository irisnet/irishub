test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@cd ${CURRENT_DIR}/e2e && go test -mod=readonly -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-nondeterminism-fast:
	@echo "Running non-determinism test..."
	@cd ${CURRENT_DIR}/e2e && go test -mod=readonly -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-import-export:
	@echo "Running application import/export simulation. This may take several minutes..."
	@cd ${CURRENT_DIR}/e2e && go test -mod=readonly -run TestAppImportExport -Enabled=true \ 
	    -NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-after-import:
	@echo "Running application simulation-after-import. This may take several minutes..."
	@cd ${CURRENT_DIR}/e2e && go test -mod=readonly -run TestAppSimulationAfterImport -Enabled=true \ 
	    -NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h