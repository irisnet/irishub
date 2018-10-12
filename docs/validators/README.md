# Validators

## Tendermint & Cosmos-SDK

Tendermint is software for securely and consistently replicating an application on many machines. Tendermint is designed to be easy-to-use, simple-to-understand, highly performant, and useful for a wide variety of distributed applications.



## What is a validator

In IRIS network , a validator is responsible for for creating new blocks and verifying transactions. The IRIS network will keep generating value when its validators could keep the whole network secure.
Validator candidates can bond their own IRIS to become a validator.


## What is a delegator

People that cannot, or do not want to run validator operations, can still participate in the staking process as delegators. Indeed, validators are not chosen based on their own stake but based on their total stake, which is the sum of their own stake and of the stake that is delegated to them. This is an important property, as it makes delegators a safeguard against validators that exhibit bad behavior. If a validator misbehaves, its delegators will move their Atoms away from it, thereby reducing its stake. Eventually, if a validator's stake falls under the top 100 addresses with highest stake, it will exit the validator set.


## States for Validator


After a validator is created with a `create-validator` transaction, it can be in the following states:

![states](../pics/states.png)

After a validator is created with a create-validator transaction, it can be in five states:

* `unbonded` & `unjailed` : 
* `bonded`:
* `unbonding`& `unjailed`:
* `unbonding`& `jailed`:
* `unbonded`&`jailed`:

Validator is in the active set and participates in consensus. Validator is earning rewards and can be slashed for misbehaviour.
unbonding: Validator is not in the active set and does not participate in consensus. Validator is not earning rewards, but can still be 
slashed for misbehaviour. This is a transition state from bonded to unbonded. If validator does not send a rebond transaction while in 
* unbonding mode, it will take three weeks for the state transition to complete.
* unbonded: Validator is not in the active set, and therefore not signing blocs. 
Validator cannot be slashed, and does not earn any reward. It is still possible to delegate Atoms to this validator. Un-delegating 
from an unbonded validator is immediate.

Once a user execute `create-validator` transaction with pubkey of a fully synced node, 
its state become `unbonded` & `unjailed`. If it's voting power is in the top 100 of all the candidates, then the state will change to `Bonded`.
But if the voting power is not enough to make to top 100, then the state will change to `unbonding`& `unjailed`. If the state keeps unchange for 3 weeks, 
then the state will change back to `unbonded` & `unjailed`. During this time, the validator could use additional delegations or add self-delegation. 
If the validator get jailed for slashing or unbond all of his self-delegation, then his state will change to `unbonding`& `jailed`. The slashing conditions
are explained below. The validator could use a `unrevoke` command to unjailed himself. His voting power will be reduced by a portion. 
If he could still remain in the top 100 validator candidates, then his state is `Bonded`, otherwise it's `unbonding`& `unjailed`.
However, if the validator doesn't unjail himself in 3 weeks, his state will be `unbonded`&`jailed`. But he could still unjail himself later. 
If he unjailed himself from the state `unbonded`&`jailed`, his state will be `Bonded` if his voting power is within top 100.

### Slashing conditions
While a validator’s goal is staying online, we will test slashing. Slashing is a punitive function that is triggered by a validator ’s bad actions. Getting slashed is losing voting power. Validators will be slashed for the actions below:

* Going offline or unable to communicate with the network

* Double sign a block

## Common Operations for Validators

* **Create Validator**

All the participants could show that they want to become a validator by sending a `create-validator` transaction, the following parameters would be necessory:

The following parameters are essential:

* `Validator's PubKey`: With flag `--pubkey`, the private key associated with PubKey is used to sign blocks. 
* `Validator's Address`: With flag `--address-validator`, the address of a participant who wants to become a validator. This is the address used to identify your validator publicly. Thhis address is used to bond, unbond, claim rewards, receive delegation and participate in governance.
* `Validator's name` : With flag `--pubkey`, default value: [do-not-modify]
* `Validator's self bond tokens`: With flag `--amount`, this value should be more than 0 and it's under unit `iris-atto`


The following parameters are optional:
* `Validator's website`: With flag `--website`
* `Validator's details`: With flag  `--details`

How to get the `PubKey` of your node?
```
iris tendermint show_validator --home={path-to-your-home}
```

The example output is the following:
```
fvp1zcjduepqcxd82mjnsnqfhwzja2d3y690ec6scw64xpg2uqkjx3rl0g0p2lwsprxnnf
```

In summary, an example of the `create-validator` command to bond 10IRIS is the following:
```
iriscli stake create-validator  --address-delegator={address1} --address-validator={address1} --name={name} --chain-id=fuxi-3001 --from=name --pubkey={pubkey} --gas=2000000 --fee=40000000000000000iris --amount=10000000000000000000iris 
```
* Edit Validator Information

Validators could edit its information.

The following parameters are optional:

* keybase signature of the validator key holder: With flag `--keybase-sig`,default value: [do-not-modify]
* the official website of the validator operator: With flag `--website`,default value: [do-not-modify]
* the name of validator: With flag `--moniker`,default value: [do-not-modify]


The command is the following:
```$xslt
iriscli stake edit-validator --address-validator={address-validator} --chain-id=fuxi-3001 --from=name  --details=details --gas=2000000 --fee=40000000000000000iris 
```

For each validator, the voting power is the sum of self-bonded token and delegated tokens. 


* View Validator Description

Anyone can view certain validator's information with the following command:
```$xslt
iriscli stake validator --address-validator={address-validator} --chain-id=fuxi-3001
```
 
* Track Validator Signing Information
In order to keep track of a validator's signatures in the past you can do so by using the signing-info command:

The command is the following:
```
iriscli stake signing-information <validator-pubkey> --chain-id=fuxi-3001
```

* Unrevoke Validator
When a validator is Revoked for downtime, you must submit an Unrevoke transaction in order to be able to get block rewards again.
In fuxi-3001, if your node missed the previous 5000 blocks in the last 10000 blocks, you will get revoked.
To query your validator's info:
```$xslt
Validator
Owner: 
Validator: 
Revoked: true
Status: Bonded
Tokens: 211.5500000000
Delegator Shares: 2800.3740578441
Description: {}
Bond Height: 642721
Proposer Reward Pool:
Commission: 0/1
Max Commission Rate: 0/1
Commission Change Rate: 0/1
Commission Change Today: 0/1
Previous Bonded Tokens: 0/1
```

The `Revoked` field of an offline validator will be `true`. 

To unrevoke your validator node, the command is the following:

```
iriscli stake unrevoke {address-validator} --from=name  --chain-id=fuxi-3001
```