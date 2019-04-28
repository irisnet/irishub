# Delegators

## What is a delegator?
People that cannot, or do not want to run validator operations, can still participate in the staking process as delegators. Indeed, validators are not chosen based on their own stake but based on their total stake, which is the sum of their own stake and of the stake that is delegated to them. This is an important property, as it makes delegators a safeguard against validators that exhibit bad behavior. If a validator misbehaves, its delegators will move their IRIS tokens away from it, thereby reducing its stake. Eventually, if a validator's stake falls under the top 100 addresses with highest stake, it will exit the validator set.

## States for a Delegator

Delegators have the same state as their validator.

Note that delegation are not necessarily bonded. Tokens of each delegator can be delegated and bonded, delegated and unbonding, delegated and unbonded, or loose. 

## Common operation for Delegators

* Delegation

To delegate 10iris token to a validator, you could run the following command:
```$xslt
iriscli stake delegate --address-delegator=<address-delegator> --address-validator=<address-validator> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --amount=10iris
```
Refer to [delegate](../cli-client/stake/delegate.md)


* Query Delegation

```$xslt
iriscli stake delegation --address-delegator=<address-delegator> --address-validator=<address-validator> --chain-id=<chain-id> 
```
Refer to [delegation](../cli-client/stake/delegation.md)


* Re-delegate 

Once a delegator has delegated his own IRIS to certain validator, he/she could change the destination of delegation at anytime. If the transaction is executed, the delegation will be placed at the other's pool after the specified period of the system parameter `unbonding_time`. 
 
```$xslt
iriscli stake redelegate --addr-validator-dest=<addr-validator-dest>  --addr-validator-source=<addr-validator> --address-delegator=<address-delegator>  --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --shares-amount=100 
```

Refer to [redelegate](../cli-client/stake/redelegate.md)

* Unbond Delegation

Once a delegator has delegated his own IRIS to certain validator, he/she could withdraw the delegation at anytime. If the transaction is executed, the delegation will become liquid after after the specified period of the system parameter `unbonding_time`.  

Unbond 50% shares：
```$xslt
iriscli stake unbond begin  --address-validator=<address-validator> --address-delegator=<address-delegator> --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --shares-percent=0.5
```

You could check that the balance of delegator has increased.

Refer to [unbond](../cli-client/stake/unbond.md)
