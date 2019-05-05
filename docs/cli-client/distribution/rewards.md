# iriscli distribution rewards

## Description

Query all the rewards of validator or delegator

## Usage

```
iriscli distribution rewards <address> <flags>
```

Print help messages:
```
iriscli distribution rewards --help
```

## Examples

```
iriscli distribution rewards <address>
```

Example response:
```json
{
  "total_rewards": [
    {
      "denom": "iris-atto",
      "amount": "235492744310548933957"
    }
  ],
  "delegations": [
    {
      "validator": "iva1q7602ujxxx0urfw7twm0uk5m7n6l9gqsgw4pqy",
      "reward": [
        {
          "denom": "iris-atto",
          "amount": "2148801198916275"
        }
      ]
    }
  ],
  "commission": [
    {
      "denom": "iris-atto",
      "amount": "235490595509350017682"
    }
  ]
}
```