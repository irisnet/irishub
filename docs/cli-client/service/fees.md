# iriscli service fees 

## Description

Query return and incoming fee of a service provider address

## Usage

```
iriscli service fees <service_provider_address>
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | help for fees                                                                                                                                         |          |

## Examples

### Query service fees

```shell
iriscli service fees <service_provider_address>
```

output:

```json
{
  "returned_fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ],
  "incoming_fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ]
}
```

