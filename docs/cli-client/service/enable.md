# iriscli service enable 

## Description

Enable an unavailable service binding

## Usage

```
iriscli service enable [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit string      |                         | [string] deposit of binding, will add to the current deposit balance                                                                                  |          |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| -h, --help            |                         | help for enable                                                                                                                                       |          |

## Examples

### Enable a unavailable service binding
```shell
iriscli service enable --chain-id=test-irishub  --from=node0 --fee=0.004iris --def-chain-id=test-irishub --service-name=test-service
```

After that, you're done with Enabling a available service binding.

```json
{
   "tags": {
     "action": "service-enable",
     "completeConsumedTxFee-iris-atto": "\"100720000000000\""
   }
 }
```