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
* **initial_height** The height to init

## Consensus Params

* **block**
  * `max_bytes` The max size of a block.
  * `max_gas`  The maximum gas quantity of a block. Its default value is -1 which means no gas limit. If the accumulation of gas consumption exceeds the block gas limit, the transaction and all subsequent transactions in the same block will fail to deliver.
  * `time_iota_ms` Minimum time increment between consecutive blocks (in milliseconds).
* **evidence**   The lifecycle of deception evidence in the block
  * `max_age_num_blocks` Max age of evidence, in blocks.
  * `max_age_duration`  Max age of evidence, in time.
  * `max_bytes` The maximum size of total evidence in bytes that can be committed in a single block.
* **validator**  The information of validator
  * `pub_key_types` The public key types validators can use.

## App State

* **auth** Params related to the system

* **bank** Params related to bank

* **capability** Params related to capability

* **coinswap** Params related to coinswap

* **crisis** Params related to crisis

* **distribution** Params related to distribution & commission

* **evidence** Params related to evidence

* **genutil** Params related to genutil

* **gov**  Params related to on-chain governance

* **guardian** Params related to guardian

* **htlc** Params related to  htlc

* **ibc** Params related to  ibc

* **nft** Params related to  nft

* **mint**  Params related to inflation

* **oracle**  Params related to oracle

* **params**  Params related to params

* **random**  Params related to random

* **record**  Params related to record

* **service**  Params related to service

* **slashing**  Params related to slashing validators

* **staking**  Params related to staking

* **token**  Params related to token

* **transfer**  Params related to transfer

* **upgrade** Params related to upgrade

* **vesting** Params related to vesting



Parameters that can be governed: [Gov Parameters](gov-params.md)

## Gentxs

Gentxs contains the transaction set of creating validators in genesis block.
The IRISnet provides robust tools for bootstrapping the identities that will start chain via the `gen-tx` process. `gen-tx` or a Genesis Transaction is cryptographically signed transactions that are executed during chain initialization that generate a starting set of validators.
The gen-txs are artifacts that prove that the holders of accounts consent in launching the network and that they put capital at risk in the process.
