# Genesis File

The Genesis file is the basis for the entire network initialization，which contains most info for creating a Genesis block (such as ChainID, consensus params,app state), initialize account balances, parameters for each module, and validators info.

## Basic State

* **genesis_time** The time to build Genesis file
* **chain_id**     Blockchain’s ID

## Consensus Params

* **block_size** 
  * `max_bytes` The max size of a block
  * `max_gas`  The maximum Gas quantity of a block. Its default value is -1 which means no gas limit. If the accumulation of comsumed gas exceeds the block gas limit, the transaction and all subsequent transactions in the same block will fail to deliver. 
* **evidence**   The lifecycle of deception evidence in the block

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

* **service**  Params related to service
  * `MaxRequestTimeout`   The max of waiting blocks for service invocation
  * `MinProviderDeposit`  The min deposit for service binding
  * `ServiceFeeTax` The service fee tax
  
* **guardian** Params related to guardian
  * `profilers` The profiler list
  * `trustees` The trustees list
  
## Gentxs

Gentxs contains the transaction set of creating validators in genesis block. 