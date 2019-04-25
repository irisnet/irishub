# iriscli service bindings

## Description

Query service bindings

## Usage

```
iriscli service bindings <flags>
```

## Flags

| Name, shorthand | Default                    | Description                                                         | Required |
| --------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --def-chain-id  |                            | the ID of the blockchain defined of the service            | Yes      |
| --service-name  |                            | service name                                               | Yes      |
| --help, -h      |                            | help for bindings                                                   |          |

## Examples

### Query service binding list

```shell
iriscli service bindings --def-chain-id=<chain-id> --service-name=<service-name>
```

After that, you will get a binding list of the service definition.

```json
[{
	"def_name": "test-service",
	"def_chain_id": "test",
	"bind_chain_id": "test",
	"provider": "iaa1ydhmma8l4m9dygsh7l08fgrwka6yczs0se0tvs",
	"binding_type": "Local",
	"deposit": [{
		"denom": "iris-atto",
		"amount": "1000000000000000000000"
	}],
	"price": [{
		"denom": "iris-atto",
		"amount": "1000000000000000000"
	}],
	"level": {
		"avg_rsp_time": "10000",
		"usable_time": "100"
	},
	"available": true,
	"disable_height": "0"
}]
```