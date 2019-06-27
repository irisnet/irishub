# iriscli asset query-gateways

## Introduction

Query all the gateways created by the specified owner

## Usage

```
iriscli asset query-gateways [flags]
```

Print help messages:
```
iriscli asset query-gateways --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --owner           | Address | false     |        | the owner address to be queried|


## Examples

```
iriscli asset query-gateways --owner=faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3
```

Output:
```txt
Gateways for owner faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3:
  Moniker: tgw, Identity: exchange, Details: testgateway, Website: http://testgateway.io
  Moniker: tgwx, Identity: exchange, Details: testgateway2, Website: http://testgateway2.io
```

```json
[
  {
    "owner": "faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3",
    "moniker": "tgw",
    "identity": "exchange",
    "details": "testgateway",
    "website": "http://testgateway.io"
  },
  {
    "owner": "faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3",
    "moniker": "tgwx",
    "identity": "exchange",
    "details": "testgateway2",
    "website": "http://testgateway2.io"
  }
]
```