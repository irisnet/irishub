# iriscli service bind 

## Description

Create a new service binding

## Usage

```
iriscli service bind <flags>
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | the average service response time in milliseconds                                                                                               |  Yes     |
| --bind-type           |                         | type of binding, valid values can be Local and Global                                                                                        |  Yes     |
| --def-chain-id        |                         | the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit             |                         | deposit of binding                                                                                                                           |  Yes     |
| --prices              |                         | prices of binding, will contains all method                                                                                                 |          |
| --service-name        |                         | service name                                                                                                                                 |  Yes     |
| --usable-time         |                         | an integer represents the number of usable service invocations per 10,000                                                                       |  Yes     |

## Examples

### Add a binding to an existing service definition
In service binding, you need to define `deposit`, the minimum mortgage amount of this `deposit` is `price` * `MinDepositMultiple` (defined by system parameters, can be modified through governance)

```shell
iriscli service bind --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<service-define-chain-id> --bind-type=Local --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999
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

