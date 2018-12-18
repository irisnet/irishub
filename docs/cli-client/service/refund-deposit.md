# iriscli service refund-deposit 

## Description

Refund all deposit from a service binding

## Usage

```
iriscli service refund-deposit [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         |  [string] service name                                                                                                                                |  Yes     |
| -h, --help            |                         |  help for refund-deposit                                                                                                                              |          |

## Examples

### Refund all deposit from an unavailable service binding
```shell
iriscli service refund-deposit --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
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