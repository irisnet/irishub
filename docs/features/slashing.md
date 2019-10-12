# Slashing

## Introduction

Collect the validator's abnormal behavior and implement the corresponding slashing mechanism according to the type of abnormal behavior.

There are three main types:

1. The validator node does not participate in the network consensus for a long time.
2. Voted multiple times on the same consensus process, and these votes contradict each other
3. The validator node disturbs the network consensus by packing illegal transactions into the block.

## Punishment mechanism

1. Calculate the number of tokens bonded to the validator node based on the voting power owned by the current validator.
2. Punish validator  with a certain percentage of the token and kick it out of the validator set; at the same time prohibit the validator from re-entering the validator set for a period of time, a process known as the jail validator.
3. For different types of abnormal behavior, different penalty proportion and jail time are used.
4. Penalty rules:

    4.1 If the total number of tokens which bonded to current validator is A and the penalty ratio is B, then the maximum number of tokens that can be punished is `A*B`.

    4.2 If there is unbonding delegation and redelegation in the unbonding period at the current height, then the tokens which in unbonding period are penalized by the ratio B.

    4.3 The total number of tokens penalized for unbonding delegation and redelegation is S. If S is less than `A*B`, the validator token punished will be `A*B-S`. Otherwise, the validator bonded token will not be penalized.

## Long Downtime

In the fixed time window `SignedBlocksWindow`, the ratio of the time of the validator's absence from the block is less than the value of `MinSignedPerWindow`, the validator's bonded token will be penalized in the `SlashFractionDowntime` ratio, and the validator will be jailed. Until the jail time exceeds `DowntimeJailDuration`, the validator can be released by executing `unjail` command.

**parameters:**

* `SignedBlocksWindow` default: 20000
* `MinSignedPerWindow` default: 0.5
* `DowntimeJailDuration` default: 2Days
* `SlashFractionDowntime` default: 0.005

## Double Sign

When executing a block, it receives evidence that a validator has voted for conflicting votes of the same round at the same height. If the time of the evidence from the current block time is less than `MaxEvidenceAge`, the validator's bonded token will be penalized in the `SlashFractionDoubleSign` ratio, and the validator will be jailed. Until the jail time exceeds `DoubleSignJailDuration`, the validator can be released by executing `unjail` command.

**parameters:**

* `MaxEvidenceAge` default: 1Day
* `DoubleSignJailDuration` default: 5Days
* `SlashFractionDoubleSign`default: 0.01

## Proposer Censorship

If the node is in the process of processing a new block, it detects if any transaction does not pass `txDecoder`, `validateTx`, `validateBasicTxMsgs`, the validator's bonded token will be slashed by `SlashFractionCensorship` percent, and the validator will be jailed. Until the jail time exceeds `CensorshipJailDuration`, the validator can be unjailed by executing the `unjail` command after jailing period.

* `txDecode` Deserialization of Tx
* `validateTx` Size limit for Tx
* `validateBasicTxMsgs` Basic check on msg in Tx

**parameters:**

* `CensorshipJailDuration` default: 7Days
* `SlashFractionCensorship` default: 0.02
