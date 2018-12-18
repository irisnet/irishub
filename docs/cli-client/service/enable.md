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
iriscli service enable --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

After that, you're done with Enabling a available service binding.

```txt
Committed at block 599 (tx hash: 87AD0FE53A665114EB3CF40505821C63952F56E9E7EF844481167C1D7B026432, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 5447,
   "codespace": "",
   "tags": {
     "action": "service_enable"
   }
 })
```