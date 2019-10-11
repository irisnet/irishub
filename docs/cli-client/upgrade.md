# iriscli upgrade

This module is used to query the software upgrade status.

## Available Commands

| Name                                            | Description                             |
| ----------------------------------------------- | --------------------------------------- |
| [info](#iriscli-upgrade-info)                   | Query the information of upgrade module |
| [query-signals](#iriscli-upgrade-query-signals) | Query the information of signals        |

## iriscli upgrade info

### Query the version the chain is running

```bash
iriscli upgrade info
```

Then it will show the current protocol information and the protocol information is currently being prepared for upgrade, e.g.

```json
{
  "version": {
    "ProposalID": "1",
    "Success": true,
    "Protocol": {
      "version": "0",
      "software": "https://github.com/irisnet/irishub/tree/v0.10.0",
      "height": "1"
    }
  },
  "upgrade_config": {
    "ProposalID": "3",
    "Definition": {
      "version": "1",
      "software": "https://github.com/irisnet/irishub/tree/v0.10.1",
      "height": "8000"
    }
  }
}
```

## iriscli upgrade query-signals

Query the current information of upgrade signals.

**Flags:**

| Name, shorthand | Default | Description        | Required |
| --------------- | ------- | ------------------ | -------- |
| --detail        | false   | Details of signals |          |

### Query the tally of upgraded voting power

```bash
iriscli upgrade query-signals
```

Example Output:

```bash
signalsVotingPower/totalVotingPower = 0.5000000000
```

### Query the detail of upgrade signals

```bash
iriscli upgrade query-signals --detail
```

Example Output:

```bash
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
siganalsVotingPower/totalVotingPower = 0.5000000000
```
