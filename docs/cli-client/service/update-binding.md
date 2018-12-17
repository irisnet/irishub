# iriscli service update-binding 

## Description

Update a service binding

## Usage

```
iriscli service update-binding [flags]
```

## Flags
| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int] the average service response time in milliseconds                                                                                               |          |
| --bind-type           |                         | [string] type of binding, valid values can be Local and Global                                                                                        |          |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit             |                         | [string] deposit of binding, will add to the current deposit balance                                                                                  |          |
| --prices              |                         | [strings] prices of binding, will contains all method                                                                                                 |          |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| --usable-time         |                         | [int] an integer represents the number of usable service invocations per 10,000                                                                       |          |
| -h, --help            |                         | help for update-binding                                                                                                                               |          |

## Examples

### Update an existing service binding
```shell
iriscli service update-binding --chain-id=test-irishub --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test-irishub --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

After that, you're done with updating an existing service binding.

```json
{
   "tags": {
     "action": "service-update-binding",
     "completeConsumedTxFee-iris-atto": "\"100840000000000\""
   }
 }
```

