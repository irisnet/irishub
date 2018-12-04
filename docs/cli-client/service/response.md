# iriscli service response 

## Description

Query a service response

## Usage

```
iriscli service response [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] the ID of the blockchain that the service invocation initiated                                                                                              |  Yes     |
| --request-id          |                         | [string] the ID of the service invocation                                                                                                                                 |  Yes     |
| -h, --help            |                         | help for response                                                                                                                                         |          |

## Examples

### Query a service response
```shell
iriscli service response --request-chain-id=test --request-id=635-535-0
```

After that, you will get the response by specified parameters.

```json
{
  "type": "iris-hub/service/SvcResponse",
  "value": {
    "req_chain_id": "test",
    "request_height": "535",
    "request_intra_tx_counter": 0,
    "expiration_height": "635",
    "provider": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "consumer": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "output": "q80=",
    "error_msg": null
  }
}
```

