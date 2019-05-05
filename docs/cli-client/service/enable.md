# iriscli service enable 

## Description

Enable an unavailable service binding

## Usage

```
iriscli service enable <flags>
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit string      |                         | deposit of binding, will add to the current deposit balance                                                                                  |          |
| --service-name        |                         | service name                                                                                                                                 |  Yes     |

## Examples

### Enable a unavailable service binding

Enable an unavailable service binding alone with 10iris additional deposit

```shell
iriscli service enable --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --def-chain-id=<service_define_chain_id> --service-name=<service_name> --deposit=10iris 
```


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