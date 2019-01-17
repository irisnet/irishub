# Changelog

## 0.10.2

*January 17th, 2019*

- [iris] The proposer must deposit 30% of the mindeposit when submitting the proposal


## 0.10.1

*January 17th, 2019*

- [iriscli] Fix issue about query validator information
- [iriscli] Fix cli query proposals error


## 0.10.0

*January 16th, 2019*

- [iris] Add flag --output-file to save export result and ensure result is consistent
- [iris] Improve invariant checking coverage and fix distribution bugs
- [iriscli] Make the result of `iriscli tendermint tx` readable
- [iriscli] Query cmd return details about software upgrade and tax usage proposal
- [tendermint] Fix the inconformity of too many evidences check
- [tendermint] Fix replay bug of `iris export`


## 0.10.0-rc0

*January 8th, 2019*

FEATURES:

- [iris] Make more validation about the `MsgCreateValidator` in CollectStdTxs
- [iris] Remove loosen token in stake pool, use bank to calculate the total loosen token
- [iris] Implement the block mint token-economics
- [iris] Add the service slash feature
- [iris] Redesign and implement the governance module to setup the new voting, tally, and penalty rules for each level of proposals
- [iris] Refactor and redefined all the gov/slashing/service/stake/distribution and gasPrice params 
- [iris] Make gov data types codec wires usable across different protocol versions
- [iris] Don't export the unfinished proposals and refund the deposits of these proposals before export snapshot 
- [iris] Refund service fee and deposit before export service state
- [iris] Add invariant checking level into makefile
- [iris] Only the genesis type profiler/trustee can initiate the addition or deletion (rather than prohibiting) transactions of the minor type profiler/trustee record. Everyone can view the profiler/trustee list
- [iris] Make sure the destination address is a trustee when the TaxUsage proposal execute
- [iris] Remove the record module
- [iris] Add `iris start --replay-last-block` to reset the app state by replay the last block
- [iris] Add `iris export --height` to export the snapshot of any block height even beyond the maximum cached historical version

- [iriscli] Add cli cmd to query the software upgrade signal status
- [iriscli] Make flag deposit not be required in the gov submit-proposal cmd
- [iriscli] Add token stats query cmd and lcd interface
- [iriscli] Replace decimal with int coins in distribution withdraw tags
- [iriscli] Add the sync tx broadcast type as the default mode in iriscli
- [iriscli] Add burn token cmd and lcd api
- [iriscli] Remove set-withdraw-addr sub-command

- [tendermint] Update tendermint to v0.27.3
- [test] Run cli test suite in parallel


BUG FIXES:

- Withdraw commission on self bond removal
- Use address instead of bond height / intratxcounter for deduplication
- Removal of mandatory self-delegation reward
- Fix bug of the tx result tags
- Fix absence proof verification
- Avoid to export account with no coin
- Correctly reset jailed-validator bond height / unbonding height on export-for-zero-height
- If a validator is jailed, distribute no reward to it
- Fix issue that miss checking the first one in Coins


## 0.9.1-patch01

*January 7th, 2019*

- Hotfix bug of software upgrade


## 0.9.1

*January 4th, 2019*

- Add cli cmd to query the software upgrade signal status
- Remove the text proposal


## 0.9.0

*December 27th, 2018*

- Refactor the gov types
- Make the deposit flag not be required in the gov submit-proposal cmd
- Add withdraw address into the withdraw tags list
- Fix the monitor bug


## 0.9.0-rc0

*December 19th, 2018*

BREAKING CHANGES:

- Use `iristool` to replace the original `irisdebug` and `irismon`
- `iris init` must specify moniker
 
FEATURES:

- [iriscli] Optimize the way tags are displayed
- [iriscli] Add `iriscli stake delegations-to [validator-addr]` and `/stake/validators/{validatorAddr}/delegations` interfaces
- [iris] Application framework code refactoring
- [iris] Add a new mechanism to distribute service fee tax
- [iris] Slashing module supports querying slashing history
- [iris] Gov module adds TxTaxUsageProposal/SoftwareHaltProposal proposals
- [iris] Export and import blockchain snapshot at any block height
- [iris] Redesigned to implement class 2 software upgrade
- [iris] Restrict the block gas limit
- [iris] Improve tx search to support multiple tags
- [iris] Improve the default behavior of iris --home
- [iris] `iris tendermint show-address` output begins with `fca`
- [iris] Restrict the number of signatures on the transaction
- [iris] Add a check for the validator private key type and reject the unsupported private key type
- [tendermint] Update tendermint to v0.27.0

BUG FIXES:

- Add chain-id value checking for sign command
- Specify the required flags for cmds `query-proposal`, `query-deposit` and `query-vote`


## 0.8.0

*December 13th, 2018*

- Upgrade tendermint to v0.27.0-dev1


## 0.8.0-rc0

*December 3rd, 2018*

BREAKING CHANGES:

- Genesis.json supports any unit format of IRIS CoinType
- The configuration information of the bech32 prefix is dynamically specified by the environment variable
- Improvement of File/directory path specification and the exception handler

FEATURES:

- Upgrade cosmos-sdk to v0.26.1-rc1 and remove the cosmos-sdk dependency
- Upgrade tendermint denpendency to v0.26.1-rc3
- View the current available withdraw balance by simulation mode
- Command line and LCD interface for service invocation request and query
- Implement guardian module for some governance proposal
- Added command add-genesis-account to configure account for genesis.json
- New proposal TerminatorProposal to terminate network consensus


## 0.7.0

*November 27th, 2018*

- Add broadcast command in bank
- Impose upgrade proposal with restrictions
- Fix bech32 prefix error in irismon
- Improve user documents

## 0.7.0-rc0

*November 19th, 2018*

BREAKING CHANGES:
* [iris] New genesis workflow
* [iris] Validator.Owner renamed to Validator. Validator operator type has now changed to sdk.ValAddress
* [iris] unsafe_reset_all, show_validator, and show_node_id have been renamed to unsafe-reset-all, show-validator, and show-node-id
* [iris]Rename "revoked" to "jailed"
* [iris]Removed CompleteUnbonding and CompleteRedelegation Msg types, and instead added unbonding/redelegation queues to endblocker
* [iris]Removed slashing for governance non-voting validators
* [iris]Validators are no longer deleted until they can no longer possibly be slashed
* [iris]Remove ibc module
* [iris]Validator set updates delayed by one block
* [iris]Drop GenesisTx in favor of a signed StdTx with only one MsgCreateValidator message

FEATURES:
* Upgrade cosmos-sdk denpendency to v0.26.0
* Upgrade tendermint denpendency to v0.26.1-rc0
* [docs]Improve docs
* [iris]Add token inflation
* [iris]Add distribution module to distribute inflation token and collected transaction fee
* [iriscli] --from can now be either an address or a key name
* [iriscli] Passing --gas=simulate triggers a simulation of the tx before the actual execution. The gas estimate obtained via the simulation will be used as gas limit in the actual execution.
* [iriscli]Add --bech to iriscli keys show and respective REST endpoint to
* [iriscli]Introduced new commission flags for validator commands create-validator and edit-validator
* [iriscli]Add commands to query validator unbondings and redelegations
* [iriscli]Add rest apis and commands for distribution

BUG FIXES:
* [iriscli]Mark --to and --amount as required flags for iriscli bank send
* [iris]Add general merkle absence proof (also for empty substores)
* [iris]Fix issue about consumed gas increasing rapidly
* [iris]Return correct Tendermint validator update set on EndBlocker by not including non previously bonded validators that have zero power
* [iris]Add commission data to MsgCreateValidator signature bytes


## 0.6.0

*November 1st, 2018*

- Use --def-chain-id flag to reference the blockchain defined of the iService
- Fix some bugs about iservice definition and record
- Add cli and lcd test for record module
- Update the user doc of iservice definition and record

## 0.6.0-rc0

*October 24th, 2018*

BREAKING CHANGES:

- [monitor] Use new executable binary in monitor

FEATURES:

- [record] Add the record module of the data certification on blockchain
- [iservice] Add the feature of iService definition
- [cli] Add the example description in the cli help
- [test] Add Cli/LCD/Sim auto-test

BUG FIXES:

- Fix software upgrade issue caused by tx fee
- Report Panic when building the lcd proof
- Fix bugs in converting validator power to byte array
- Fix panic bug in wrong account number


## 0.5.0-rc1

*October 11th, 2018*

FEATURES:

- Make all the gov and upgrade parameters can be configured in the genesis.json

BUG FIXES

- Add check for iavl proof and value before building multistore proof


## 0.5.0-rc0

*September 30th, 2018*

BREAKING CHANGES:

- [cointype] Introduce the cointype of iris:
  - 1 iris = 10^18 iris-atto
  - 1 iris-milli = 10^15 iris-atto
  - 1 iris-micro = 10^12 iris-atto
  - 1 iris-nano = 10^9 iris-atto
  - 1 iris-pico = 10^6 iris-atto
  - 1 iris-femto = 10^3 iris-atto

FEATURES:

- [tendermint] Upgrade to Tendermint v0.23.1-rc0
- [cosmos-sdk] Upgrade to cosmos-sdk v0.24.2
    - Move the previous irisnet changeset about cosmos-sdk into irishub
- [irisdebug] Add irisdebug tool
- [LCD/cli] Add the proof verification to the LCD and CLI
- [iparam] Support the modification of governance parameters of complex data type through governance, and the submission of modified proposals through json config files
- [software-upgrade] Software upgrade solutions of the irisnet


## 0.4.2

*September 22th, 2018*

BUG FIXES

- Fix consensus failure due to the double sign evidence be broadcasted before the genesis block

## 0.4.1

*September 12th, 2018*

BUG FIXES

- Missing to set validator intraTxCount in stake genesis init


## 0.4.0

*September 6th, 2018*

BREAKING CHANGES:

- [cosmos-sdk] Upgrade to cosmos-sdk v0.23.0
    - Change the address prefix format:
      - cosmosaccaddr --> faa
      - cosmosaccpub --> fap
      - cosmosvaladdr --> fva
      - cosmosvalpub --> fvp
    - Adjust the Route & rootMultiStore Commit for software upgrade
    - Must specify gas and fee in both command line and rest api
    - The fee should be iris token and the token amount should be no less than 2*(10^10)*gas

FEATURES:

- [tendermint] Upgrade to Tendermint v0.22.6
    - Store the pre-state to support the replay function
- [cosmos-sdk] Upgrade to cosmos-sdk v0.23.0
    - Add the paramProposal and softwareUpgradeProposal in gov module
    - Improve fee token mechanism to more reasonably deduct transaction fee and achieve more ability to defent DDOS attack.
    - Introduce the global parameter module

BUG FIXES

- Default account balance in genesis
- Fix iris version issue
- Fix the unit conflict issue in slashing
- Check the voting power when create validator
- Fix evidence amimo register issue


## 0.4.0-rc2

*Sep 5th, 2018*

BUG FIXES

- Fix evidence amimo register issue


## 0.4.0-rc1

*Aug 27th, 2018*

BUG FIXES

- Default account balance in genesis
- iris version issue
- Fix the unit conflict issue in slashing
- Check the voting power when create validator


## 0.3.0

*July 30th, 2018*

BREAKING CHANGES:

- [tendermint] Upgrade to Tendermint v0.22.2
    - Default ports changed from 466xx to 266xx
    - ED25519 addresses are the first 20-bytes of the SHA256 of the raw 32-byte pubkey (Instead of RIPEMD160)
- [cosmos-sdk] Upgrade to cosmos-sdk v0.22.0
- [monitor] Move `iriscli monitor` subcommand to `iris monitor`

FEATURES:

- [lcd] /tx/send is now the only endpoint for posing transaction to irishub; aminofied all transaction messages 
- [monitor] Improve the metrics for iris-monitor 

BUG FIXES

- [cli] solve the issue of iriscli stake sign-info

##

## 0.2.0

*July 19th, 2018*

BREAKING CHANGES:

- [tendermint] Upgrade to Tendermint v0.21.0
- [cosmos-sdk] Upgrade to cosmos-sdk v0.19.1-rc1

FEATURES:

- [lcd] code refactor

- [cli] improve sendingand querying the  transactions 

- [monitor]export data which is collected by Prometheus Server

  ​

##  