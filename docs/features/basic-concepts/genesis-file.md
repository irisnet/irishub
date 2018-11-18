# Genesis File

The Genesis file is the basis for the entire network initialization，which contains most info for creating a Genesis block (such as ChainID, consensus params,app state), initialize account balances, parameters for each module, and validators info.

## Basic State

* genesis_time The time to build Genesis file
* chain_id     Blockchain’s ID

## Consensus Params

* block_size Block size and config params of the number of Gas in the block
* evidence   The lifecycle of deception evidence in the block

## App State

* **accounts** Initialization account info

* **stake** Params related to the staking consensus
  * `loose_tokens`   The sum of unbonded tokens in the entire network
  * `unbonding_time` The time between the moment a validator begin to unbond until the moment it is unbonded successfully
  * `max_validators` The max of validators
  
* **mint**  Params related to inflation
  * `inflation_max` The max of inflation rate
  * `inflation_min` The min of inflation rate
  
* **distribution** Params related to distribution & commission

* **gov**  Params related to on-chain governance
  * `DepositProcedure`  Params in deposit period
  * `VotingProcedure`   Params in voting period
  * `TallyingProcedure` Params in tallying period

* **upgrade** Params related to upgrade
  * `switch_period` After upgrade, a switch message needs to be sent in switch_perid

* **slashing** Params related to slashing validators

* **service**  Params related to Service
  * `MaxRequestTimeout`   The max of waiting blocks for service invocation
  * `MinProviderDeposit`  The min deposit for service binding
  
## Gentxs

Gentxs contains the transaction set of creating validators in genesis block. 