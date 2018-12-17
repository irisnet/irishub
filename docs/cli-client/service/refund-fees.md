# iriscli service refund-fees 

## Description

Refund all fees from service return fees

## Usage

```
iriscli service refund-fees [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | help for refund-fees                                                                                                                                         |          |

## Examples

### Refund fees from service return fees 
```shell
iriscli service refund-fees --chain-id=test-irishub --from=node0 --fee=0.004iris
```

After that, you're done with refunding fees from service return fees.

```json
{
   "tags": {
     "action": "service-refund-fees",
     "completeConsumedTxFee-iris-atto": "\"679600000000000\""
   }
 }
```

