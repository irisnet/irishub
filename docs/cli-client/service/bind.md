# iriscli service bind 

## Description

Create a new service binding

## Usage

```
iriscli service bind [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int] the average service response time in milliseconds                                                                                               |  Yes     |
| --bind-type           |                         | [string] type of binding, valid values can be Local and Global                                                                                        |  Yes     |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit             |                         | [string] deposit of binding                                                                                                                           |  Yes     |
| --prices              |                         | [strings] prices of binding, will contains all method                                                                                                 |          |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| --usable-time         |                         | [int] an integer represents the number of usable service invocations per 10,000                                                                       |  Yes     |
| -h, --help            |                         | help for bind                                                                                                                                         |          |

## Examples

### Add a binding to an existing service definition
```shell
iriscli service bind --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

After that, you're done with adding a binding to an existing service definition.

```txt
Committed at block 567 (tx hash: A48DBD217CBB843E72CC47B40F90CE7DEEEDD6437C86A74A2976ADC9F449A034, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 5893,
   "codespace": "",
   "tags": {
     "action": "service_bind"
   }
 })
```

