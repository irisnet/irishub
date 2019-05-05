# iriscli service call 

## Description

Call a service method

## Usage

```
iriscli service call <flags>
```

## Flags

| Name, shorthand       | Default                 | Description                                                  | Required |
| --------------------- | ----------------------- | ------------------------------------------------------------ | -------- |
| --def-chain-id        |                         | the ID of the blockchain defined of the service     |  Yes     |
| --service-name        |                         | service name                                        |  Yes     |
| --method-id           |                         | the method id called                                   |  Yes     |
| --bind-chain-id       |                         | the ID of the blockchain bond of the service        |  Yes     |
| --provider            |                         | bech32 encoded account created the service binding  |  Yes     |
| --service-fee         |                         | fee to pay for a service invocation                 |          |
| --request-data        |                         | hex encoded request data of a service invocation    |          |

## Examples

### Initiate a service invocation request 

```shell
iriscli service call --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --def-chain-id=<service_define_chain_id> --service-name=<service_name> --method-id=1 --bind-chain-id=<service_bind_chain_id> --provider=<provider_address> --service-fee=1iris --request-data=<request-data>
```

After that, you're done with initiating a service invocation request.

```txt
Committed at block 54 (tx hash: F972ACA7DF74A6C076DFB01E7DD49D8694BF5AA1BA25A1F1B875113DFC8857C3, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7614,
   "codespace": "",
   "tags": {
     "action": "service_call",
     "consumer": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "provider": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "request-id": "64-54-0"
   }
 })
```

