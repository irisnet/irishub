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
iriscli service update-binding --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

After that, you're done with updating an existing service binding.

```txt
Committed at block 579 (tx hash: D95E002AF467A7C4E7F298664E8C1951522B4CB61D26B01AC9705703E75557AB, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 5449,
   "codespace": "",
   "tags": {
     "action": "service_binding_update"
   }
 })
```

