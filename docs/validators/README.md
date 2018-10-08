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

![states](../pics/states.jpg)
  
