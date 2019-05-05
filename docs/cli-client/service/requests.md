# iriscli service requests 

## Description

Query service requests

## Usage

```
iriscli service requests <flags>
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         | service name                                                                                                                                 |  Yes     |
| --bind-chain-id       |                         | the ID of the blockchain bond of the service                                                                                                                                 |  Yes     |
| --provider            |                         | bech32 encoded account created the service binding                                                                       |  Yes     |

## Examples

### Query service request list

```shell
iriscli service requests --def-chain-id=<service_define_chain_id> --service-name=<service-name> --bind-chain-id=<service_bind_chain_id> --provider=<provider_address>
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
    "provider": "iaa1f02ext9duk7h3rx9zm7av0pnlegxve8npm2k6m",
    "consumer": "iaa1f02ext9duk7h3rx9zm7av0pnlegxve8npm2k6m",
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

