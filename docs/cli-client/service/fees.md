# iriscli service fees 

## Description

Query return and incoming fee of a particular address

## Usage

```
iriscli service fees [account address]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | help for fees                                                                                                                                         |          |

## Examples

### Query service fees
```shell
iriscli service fees iaa1f02ext9duk7h3rx9zm7av0pnlegxve8npm2k6m
```

After that, you will get the service fees by specified address.

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

