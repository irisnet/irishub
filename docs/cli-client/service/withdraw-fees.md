# iriscli service withdraw-fees 

## Description

Withdraw all fees from service incoming fees

## Usage

```
iriscli service withdraw-fees [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | help for withdraw-fees                                                                                                                                         |          |

## Examples

### Withdraw fees from service incoming fees 
```shell
iriscli service withdraw-fees --chain-id=test-irishub --from=node0 --fee=0.004iris
```

After that, you're done with withdraw fees from service incoming fees.

```json
{
   "tags": {
     "action": "service-withdraw-fees",
     "completeConsumedTxFee-iris-atto": "\"679600000000000\""
   }
 }
```

