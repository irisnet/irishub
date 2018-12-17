# Software Upgrade User Guide

## Basic Function Description

The module supports the infrastructure of the blockchain software upgrade. It will be upgraded to the new version through voting at UpgradeProposal  and is fully compatible with the historical data on the blockchain.

## Interaction Process

###  Governance process of software upgrade proposal
1. Submit a software upgrade proposal
2. More details about governance process is in GOV [User Guide](governance.md)

### The process of software upgrade   
1. Install a new software.
2. Once reach the limited time, it will be counted whether the proportion of voting power of upgraded software exceeds 95%.
3. If it exceeds 95%, the software will be upgraded, otherwise the upgrade fails.
4. The validators who didn't upgrade in time need to re-download the new software and blocks synchronized.

## Usage Scenarios

### Create an environment

```
rm -rf iris                                                                         
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=upgrade-test -o --home=iris
iris start --home=iris
```

### Submit a software upgrade proposal

```
# Send an upgrade proposal
iriscli gov submit-proposal --title=Upgrade --description="SoftwareUpgrade" --type="SoftwareUpgrade" --deposit=10iris --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000 --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node
```

### Upgrade software

* Scenario 1

Implement following operations at the certain height（80 block height）：

```
# 1. Download the new version:iris1

# 2. Close the old one
kill -f iris

# 3. Install the new version，iris1 and start it（copy to bin）
iris1 start --home=iris

# 4. Upgrade automatically when reach the preset time

# 5. Query whether the current version has been successfully upgraded
iriscli upgrade info --trust-node
```

* Scenario 2

The operations in Scenario 1 haven't been implemented at the certain time (80 block height), report errors after the new version become valid:

```
# 1. Download the new version, iris1

# 2. Close the old one
kill -f iris

# 3. Install the new version iris1 and start it 

iris1 start --home=iris

# 4. Query whether the current version has been successfully upgraded
iriscli upgrade info --trust-node
```

## Command details

```
iriscli gov submit-proposal --title=Upgrade --description="SoftwareUpgrade" --type="SoftwareUpgrade" --deposit=10iris --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000 --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80
```

* `--type`  "SoftwareUpgrade" The type of Software upgrade proposals
* `--version`  The version of the new protocol
* `--software`  The software of the new protocol
* `--switch-height` The switchheight of the new protocol
* Other parameters can be referrenced in [Gov User Guide](governance.md)

```
iriscli upgrade submit-switch --name=x --from=$VADDR --proposalID=1 --chain-id=upgrade-test --fee=0.05iris --gas=20000
```

* `--proposalID` The ID of passed software upgrade proposals 

```
iris start --replay
```

* Resynchronize the block, clean the dirty AppHash

```
iriscli upgrade info --trust-node
```

* Query the version details of current software 
