# iriscli service requests 

## Description

Query service requests

## Usage

```
iriscli service requests [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| --bind-chain-id       |                         | [string] the ID of the blockchain bond of the service                                                                                                                                 |  Yes     |
| --provider            |                         | [string] bech32 encoded account created the service binding                                                                       |  Yes     |

## Examples

### Query service request list
```shell
iriscli service requests --def-chain-id=<chain-id> --service-name=test-service --bind-chain-id=test --provider=iaa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x
```

After that, you will get the active request list of the specified provider.

```json
[
  {
    "def_chain_id": "chain-jsmJQQ",
    "def_name": "test-service",
    "bind_chain_id": "chain-jsmJQQ",
    "req_chain_id": "chain-jsmJQQ",
    "method_id": 1,
    "provider": "iaa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "consumer": "iaa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "input": "Q0NV",
    "service_fee": [
      {
        "denom": "iris-atto",
        "amount": "10000000000000000"
      }
    ],
    "profiling": false,
    "request_height": "456",
    "request_intra_tx_counter": 0,
    "expiration_height": "556"
  }
]
```

