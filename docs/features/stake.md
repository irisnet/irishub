# Stake User Guide

## Introduction

This specification briefly introduces the functionality of stake module and what user should do with the provided commands.

## Core Concept

1. Voting power

	Voting power is a consensus concept. IRISHUB is a Byzantine-fault-tolerant POS blockchain network. During the consensus process, a set of validators will vote the proposal block. If a validator thinks the proposal block is valid, it will vote `yes`, otherwise, it will vote nil. The votes from different validator don't have the same weight. The weight of a vote is called the voting power of the corresponding validator.
	
2. Validator

    Validator is a full IRISHUB node. If its voting power is zero, it is just a normal full node or a validator candidate. Once its voting power is positive, then it is a real validator.
     
3. Delegator && Delegation

	People that cannot, or do not want to run validator nodes, can still participate in the staking process as delegators. After delegating some tokens to validators, delegators will gain shares from corresponding validators. Delegating tokens is also called bonding tokens to validators. Later we will have detailed description on it. Besides, a validator operator is also a delegator. Usually, a validator operator only has delegation on its own validator. But it can also have delegation on other validators.
	
4. Validator Candidates
 
	The quantity of validators can't increase without limit. Too many validators may result in low efficient consensus which slows down the blockchain TPS. So Byzantine-fault-tolerant POS blockchain network will have a limiation to the validator quantity. Usually, the value is 100. If more than 100 full nodes apply to join validator set. Then only these nodes with top 100 most bounded tokens will be real validators. Others will be validator candidates and will be descending sorted according to their bonded token amount. Once the one or more validators are kicked out from validator set, then the top candidates will be added into validator set automatically.
	
5. Bond && Unbond && Unbonding Period

	Validator operators must bond their liquid tokens to their validators. The validator voting power is proportional to the bonded tokens including both self-bonded tokens and tokens from other delegators. Validator operators can lower their own bonded tokens by sending unbond transactions. Delegators can also lower bonded token by sending unbond transactions. However, these unbonded token won't become liquid tokens immediately. After the unbond transactions are executed, the corresponding validator operators or delegators can't sending unbond transactions on the same validators again until the unbonding period is end. Usually the unbonding period is three weeks. Once the unbonding period is end, the unbonded token will become liquid token automatically. The unbonding period mechanism makes great contribution to the security of POS blockchain network. Besides, if the self-bonded token equals to zero, then the corresponding validator will be removed out of validator set.
	 
6. Redelegate

	Delegators can transfer their delegation from one validator to another one. Redelegation can be devided into two steps: ubond from first validator and bond to another validator. As we have talked above, ubond operation can't be completed immediately until unbonding period is end, which means delegators can't send another redelegation transactions immediately.
	
7. Evidence && Slash

	The Byzantine-fault-tolerant POS blockchain network assume that the Byzantine nodes possess less than 1/3 of total voting power. These Byzantine nodes must be punished. So it is necessary to collect the evidence of Byzantine behavior. According to the evidence, stake module will aotumatically slash a certain mount of token from corresponding validators and delegators. The slashed tokens are just burned. Besides, the Byzantine validators will be removed from the validator set and put into jail, which means their voting power is zero. During the jail period, these nodes are not event validator candidates . Once the jail period is end, they can send transactions to unjail themselves and become validator candidates again.
	
8. Rewards

	As a delegator, the more bonded tokens it has on validator, the more rewards it will earn. For a validator operator, it will have extra rewards: validator commission. The rewards comes from token inflation and transaction fee. As for how to calculate the rewards and how to get the rewards, please refer to [mint](mint.md) and [distribution](distribution.md).
	
## What users can do

1. Create a full node

	Please refer to [full node](../get-started/Full-Node.md) to create a full node.

2. Apply to be validator

	Firstly, you must have a wallet which has a certain amount of iris tokens. Here we assume you have import your wallet to iriscli key store. 

	Then just send a create-validator transaction. This is an example command.
	```
	iriscli stake create-validator --amount=100iris --pubkey=$(iris tendermint show-validator) --moniker=<validator name> --fee=0.004iris --chain-id=<chain-id> --from=<key name> --commission-max-change-rate=0.01 --commission-max-rate=0.2 --commission-rate=0.1
	```
	The more tokens specified by `--amount`, the more probability your full node will be a real validator. Otherwise, it will just be validator candidate.

3. Query your own validator
	
	Users can query their own validators by their wallet address. But firstly users have to convert their wallet addresses to validator operator address pattern:
	```
	iriscli keys show [key name] --bech=val
	```
	Example response:
	```
	NAME:   TYPE:   ADDRESS:                                      PUBKEY:
	faucet  local   fva1ljemm0yznz58qxxs8xyak7fashcfxf5l9pe40u    fvp1addwnpepqtdme789cpm8zww058ndlhzpwst3s0mxnhdhu5uyps0wjucaufha605ek3w
	```
	Then, example command to query validator:
	```
	iriscli stake validator fva1ljemm0yznz58qxxs8xyak7fashcfxf5l9pe40u
	```
	Example response:
	```text
    Validator 
    Operator Address: fva1ljemm0yznz58qxxs8xyak7fashcfxf5l9pe40u
    Validator Consensus Pubkey: fvp1zcjduepq8fw9p4zfrl5fknrdd9tc2l24jnqel6waxlugn66y66dxasmeuzhsxl6m5e
    Jailed: false
    Status: Bonded
    Tokens: 100.0000000000
    Delegator Shares: 100.0000000000
    Description: {node2   }
    Bond Height: 0
    Unbonding Height: 0
    Minimum Unbonding Time: 1970-01-01 00:00:00 +0000 UTC
    Commission: {{0.1000000000 0.2000000000 0.0100000000 0001-01-01 00:00:00 +0000 UTC}}
    ```
	
4. Edit validator

	```
	iriscli stake edit-validator --from=<key name> --chain-id=<chain-id> --fee=0.004iris --commission-rate=0.15 --moniker=<new name>
	```
	
5. Increase self-delegation

	```
	iriscli stake delegate --address-validator=<self-address-validator> --chain-id=<chain-id> --from=<key name> --fee=0.004iris  --amount=100iris 
	```

6. Delegate tokens to other validators

	If you just want to be a delegator, you can skip the above steps.
	```
	iriscli stake delegate --address-validator=<other-address-validator> --chain-id=<chain-id> --from=<key name> --fee=0.004iris  --amount=100iris 
	```

7. Unbond tokens from a validator

	Unbond half of total bonded token on a given validator
	```
	iriscli stake unbond --address-validator=<address-validator> --chain-id=<chain-id> --from=<key name> --fee=0.004iris  --amount=100iris --share-percent=0.5
	```

8. Redelegate tokens to another validator

	Redelegate half of total bonded token on a given validator to another one
	```
	iriscli stake redelegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --address-validator-source=<source validator address> --address-validator-dest=<destination validator address> --shares-percent=0.5
	```

For other query stake state commands, please refer to [stake cli client](../cli-client/stake/README.md)
