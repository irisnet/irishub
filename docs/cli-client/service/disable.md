# iriscli service disable 

## Description

Disable a available service binding

## Usage

```
iriscli service disable [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| -h, --help            |                         | help for disable                                                                                                                                      |          |

## Examples

### Disable a available service binding
```shell
iriscli service disable --chain-id=test-irishub  --from=node0 --fee=0.004iris --def-chain-id=test-irishub --service-name=test-service
```

After that, you're done with disabling a available service binding.

```json
{
   "tags": {
     "action": "service-disable",
     "completeConsumedTxFee-iris-atto": "\"70200000000000\""
   }
 }
```