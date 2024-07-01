CONTRACTS_DIR := $(CURDIR)/modules/token/contracts
COMPILED_DIR := $(CONTRACTS_DIR)/compiled_contracts
NODE_MODULES := $(CONTRACTS_DIR)/node_modules

# Compile and format solidity contracts for the erc20 module. Also install
# openzeppeling as the contracts are build on top of openzeppelin templates.
contracts-compile: contracts-clean dep-install create-contracts-abi

# Install openzeppelin solidity contracts
dep-install:
	@echo "Importing openzeppelin contracts..."
	@cd $(CONTRACTS_DIR) && npm install
	 
# Clean tmp files
contracts-clean:
	@rm -rf $(NODE_MODULES)

# Compile, filter out and format contracts into the following format.
create-contracts-abi:
	solc --combined-json abi,bin --optimize --optimize-runs 200 --evm-version paris --include-path $(NODE_MODULES) --base-path $(CONTRACTS_DIR)/  $(CONTRACTS_DIR)/Token.sol | jq '.contracts["Token.sol:Token"]' > $(COMPILED_DIR)/Token.json \
    && solc --combined-json abi,bin --optimize --optimize-runs 200 --evm-version paris --include-path $(NODE_MODULES) --base-path $(CONTRACTS_DIR)/  $(CONTRACTS_DIR)/TokenProxy.sol | jq '.contracts["TokenProxy.sol:TokenProxy"]' > $(COMPILED_DIR)/TokenProxy.json \
	&& solc --combined-json abi,bin --optimize --optimize-runs 200 --evm-version paris --include-path $(NODE_MODULES) --base-path $(CONTRACTS_DIR)/  $(CONTRACTS_DIR)/UpgradeableBeacon.sol | jq '.contracts["UpgradeableBeacon.sol:UpgradeableBeacon"]' > $(COMPILED_DIR)/UpgradeableBeacon.json \