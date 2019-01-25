# Slash

## Basic Function

Collect the validator's abnormal behavior and implement the corresponding slashing mechanism according to the type of abnormal behavior.

There are three main types:

1. The validator node does not participate in the network consensus for a long time.
2. Voted multiple times on the same consensus process, and these votes contradict each other
3. The validator node disturbs the network consensus by packing illegal transactions into the block.

## Punishment mechanism

1. Calculate the number of tokens bound to the validator node based on the voting power owned by the current validator.
2. Punish validator  with a certain percentage of the token and kick it out of the validator set; at the same time prohibit the validator from re-entering the validator set for a period of time, a process known as the jail validator.
3. For different types of abnormal behavior, use different penalty ratios and jail time.
4. Penalty rules:
1. If the current validator token total is A and the penalty ratio is B, then the number of tokens that the validator can punish at most is A*B.
2. If there is unbonding delegation and redelegation in the unbonding period at the current height, and the creation height of unbonding delegation and redegetation is less than the execution height of the abnormal behavior, then the tokens of the two parts are penalized by  the ratio B.
3. The total number of tokens penalized for unbonding delegation and redelegation is S. If S is less than A*B, the validator token punished will be `A*B-S`. Otherwise, the validator bonded token is not penalized.

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

