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
iriscli service disable --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

After that, you're done with disabling a available service binding.

```txt
Committed at block 588 (tx hash: 12333DCD222EB3620A2177E9A8E0D84E248BAE0D3BC445274E09A19096794A46, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3123,
   "codespace": "",
   "tags": {
     "action": "service_disable"
   }
 })
```