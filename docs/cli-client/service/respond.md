# iriscli service respond 

## Description

Respond a service method invocation

## Usage

```
iriscli service respond <flags>
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | the ID of the blockchain that the service invocation initiated                                                                                              |  Yes     |
| --request-id          |                         | the ID of the service invocation                                                                                                                                |  Yes     |
| --response-data       |                         | hex encoded response data of a service invocation                                                                       |         |

## Examples

### Respond to a service invocation 

```shell
iriscli service respond --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --request-chain-id=<request_chain_id> --request-id=<request-id> --response-data=<response-data>
```
> You can figure out the `request-id` in the return of [service call](call.md)

After that, you're done with responding to a service invocation.

```txt
Committed at block 71 (tx hash: C02BC5F4D6E74ED13D8D5A31F040B0FED0D3805AF1C546544A112DB2EFF3D9D5, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3784,
   "codespace": "",
   "tags": {
     "action": "service_respond",
     "consumer": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "provider": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "request-id": "78-68-0"
   }
 })
```

