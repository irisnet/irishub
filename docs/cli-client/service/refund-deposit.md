# iriscli service refund-deposit 

## Description

Refund all deposit from a service binding

## Usage

```
iriscli service refund-deposit <flags>
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         | service name                                                                                                                                |  Yes     |

## Examples

### Refund all deposit from an unavailable service binding

Before refunding , you should [disable](disable.md) the service binding first.

```shell
iriscli service refund-deposit --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --def-chain-id=<service_define_chain_id> --service-name=<service-name>
```

After that, you're done with refunding all deposit from a service binding.

```txt
Committed at block 17 (tx hash: 6C878E864772DE2F29725B743A8B9D40A75B41688F16C278634674653BFD1DFA, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 4787,
   "codespace": "",
   "tags": {
     "action": "service_refund_deposit"
   }
 })
```