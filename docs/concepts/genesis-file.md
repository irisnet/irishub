---
order: 4
---

# Genesis File

The Genesis file (~/.iris/config/genesis.json) is the basis for the entire network initialization, which contains most info for creating a Genesis block (such as ChainID, consensus params, app state, initial account balances, parameters for each module, and validators info).
The genesis file sets the initial parameters of any new IRIS network. Establishing a robust social consensus over the genesis file is critical to start a network.

Each genesis state starts with a list of account balances. Social consensus on these account balances must be bootstrapped from some external process be it events on another blockchain to a token generation event.

## Basic State

* **genesis_time** The time to launch
* **chain_id**     Blockchain’s ID

## Consensus Params

* **block_size**
  * `max_bytes` The max size of a block
  * `max_gas`  The maximum gas quantity of a block. Its default value is -1 which means no gas limit. If the accumulation of gas consumption exceeds the block gas limit, the transaction and all subsequent transactions in the same block will fail to deliver.
* **evidence**   The lifecycle of deception evidence in the block

## App State

* **accounts** Initial account info

* **auth** Params related to the system

* **stake** Params related to the staking consensus
  
* **mint**  Params related to inflation
  
* **distribution** Params related to distribution & commission
  
* **gov**  Params related to on-chain governance

* **upgrade** Params related to upgrade

* **slashing** Params related to slashing validators

* **asset**  Params related to asset

* **service**  Params related to service

* **guardian** Params related to guardian

Parameters that can be governed: [Gov Parameters](gov-params.md)

## Gentxs

Gentxs contains the transaction set of creating validators in genesis block.
The IRISnet provides robust tools for bootstrapping the identities that will start chain via the `gen-tx` process. `gen-tx` or a Genesis Transaction is cryptographically signed transactions that are executed during chain initialization that generate a starting set of validators.
The gen-txs are artifacts that prove that the holders of accounts consent in launching the network and that they put capital at risk in the process.
