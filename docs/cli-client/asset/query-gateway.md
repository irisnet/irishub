# iriscli asset query-gateway

## Introduction

Query a gateway which is identified by the specified moniker

## Usage

```
iriscli asset query-gateway [flags]
```

Print help messages:
```
iriscli asset query-gateway --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | true     | ""       | the unique name with a size between 3 and 8 letters|


## Examples

```
iriscli asset query-gateway --moniker tgw
```

Output:
```txt
Gateway:
  Owner:             faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3
  Moniker:           tgw
  Identity:          exchange
  Details:           testgateway
  Website:           http://testgateway.io
```

```json
{
  "type": "irishub/asset/Gateway",
  "value": {
    "owner": "faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3",
    "moniker": "tgw",
    "identity": "exchange",
    "details": "testgateway",
    "website": "http://testgateway.io"
  }
}
```