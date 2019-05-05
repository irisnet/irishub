# iriscli service binding

## Description

Query service binding

## Usage

```
iriscli service binding <flags>
```

## Flags

| Name, shorthand | Default                    | Description                                                         | Required |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --bind-chain-id |                            | the ID of the blockchain bond of the service               | Yes      |
| --def-chain-id  |                            | the ID of the blockchain defined of the service            | Yes      |
| --provider      |                            | bech32 encoded account created the service binding         | Yes      |
| --service-name  |                            | service name                                               | Yes      |
| --help, -h      |                            | help for binding                                                    |          |

## Examples

### Query a service binding

```shell
iriscli service binding --def-chain-id=<service_define_chain-id> --service-name=<service-name> --bind-chain-id=<service_bind_chain-id> --provider=<provider_address>
```

After that, you will get detail info for the service binding.

```json
{
  "type": "iris-hub/service/SvcBinding",
  "value": {
    "def_name": "test-service",
    "def_chain_id": "test",
    "bind_chain_id": "test",
    "provider": "iaa1ydhmma8l4m9dygsh7l08fgrwka6yczs0se0tvs",
    "binding_type": "Local",
    "deposit": [
      {
        "denom": "iris-atto",
        "amount": "1000000000000000000000"
      }
    ],
    "price": [
      {
        "denom": "iris-atto",
        "amount": "1000000000000000000"
      }
    ],
    "level": {
      "avg_rsp_time": "10000",
      "usable_time": "100"
    },
    "available": true,
    "disable_height": "0"
  }
}
```