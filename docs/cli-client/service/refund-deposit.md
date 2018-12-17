# iriscli service refund-deposit 

## Description

Refund all deposit from a service binding

## Usage

```
iriscli service refund-deposit [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         |  [string] service name                                                                                                                                |  Yes     |
| -h, --help            |                         |  help for refund-deposit                                                                                                                              |          |

## Examples

### Refund all deposit from an unavailable service binding
```shell
iriscli service refund-deposit --chain-id=test-irishub  --from=node0 --fee=0.004iris --def-chain-id=test-irishub --service-name=test-service
```

After that, you're done with refunding all deposit from a service binding.

```json
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "\"92280000000000\""
   }
 }
```