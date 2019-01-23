# Slash

## Basic Function

By punishing some byzantine validators, we can better maintain the healthy growth of the network and maintain the network's activeness.

There are three types of punishment:

1. Punish the validator for a long time not online
2. Punish the validator's double sign behavior
3. Punish the validator hack code to packing the garbage data into blockchain with minimal cost.

## Long Downtime

In the fixed time window `SignedBlocksWindow`, the ratio of the time of the validator's absence from the block is less than the value of `MinSignedPerWindow`，the validator's bonded token is penalized in the `SlashFractionDowntime` ratio, and the validator is jailed. Until the jail time exceeds `DowntimeJailDuratio`, the validator can be released by the unjail command.

### parameters

* `SignedBlocksWindow` default: 20000
* `MinSignedPerWindow` default: 0.5
* `DowntimeJailDuration` default: 2Days
* `SlashFractionDowntime` default: 0.005

## Double Sign

When executing a block, it receives evidence that a validator has signed the different blocks of the same round at the same height. If the time of the evidence from the current block time is less than `MaxEvidenceAge`，the validator's bonded token is penalized in the `SlashFractionDoubleSign` ratio, and the validator is jailed. Until the jail time exceeds `DoubleSignJailDuration`, the validator can be released by the unjail command.

### parameters

* `MaxEvidenceAge` default: 1Day
* `DoubleSignJailDuration` default: 5Days
* `SlashFractionDoubleSign`default: 0.01

## Propoer Censorship

If the node is in the process of executing the block, it detects that the transaction does not pass `txDecoder`, `validateTx`, `validateBasicTxMsgs`, the validator's bonded token is penalized in the `SlashFractionCensorship` ratio, and the validator is jailed. Until the jail time exceeds `CensorshipJailDuration`, the validator can be released by the unjail command.

* `txDecode` Deserialization of Tx
* `validateTx` Size limit for Tx
* `validateBasicTxMsgs` Basic check on msg in Tx

### parameters

* `CensorshipJailDuration` default: 7Days
* `SlashFractionCensorship` default: 0.02

## slash command

### unjail

If the validator was jailed and the jail time passed, release the validator by unjail command.

```
iriscli stake unjail --from=<key name> --fee=0.004iris --chain-id=<chain-id>
```
