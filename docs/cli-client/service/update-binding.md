# iriscli service update-binding 

## Description

Update a service binding

## Usage

```
iriscli service update-binding <flags>
```

## Flags
| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | the average service response time in milliseconds                                                                                               |          |
| --bind-type           |                         | type of binding, valid values can be Local and Global                                                                                        |          |
| --def-chain-id        |                         | the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit             |                         | deposit of binding, will add to the current deposit balance                                                                                  |          |
| --prices              |                         | prices of binding, will contains all method                                                                                                 |          |
| --service-name        |                         | service name                                                                                                                                 |  Yes     |
| --usable-time         |                         | an integer represents the number of usable service invocations per 10,000                                                                       |          |

## Examples

### Update an existing service binding

Update service binding alone with 10iris additional deposit

```shell
iriscli service update-binding --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --service-name=<service_name> --def-chain-id=<service_define_chain_id> --bind-type=Local --deposit=10iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999
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

