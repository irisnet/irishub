# iriscli upgrade query-signals

## Description

Query the current information of signals.

## Flags

| Name, shorthand | Default               | Description                                                                                                                                 | Required |
| --------------- | --------------------- | --------------------------------------------------------------------------------------------------------------------------------------------| -------- |
| --detail        | false                 | [bool] details of signals                                                                                                                 |          |

## Usage

```
iriscli upgrade query-signals
```

Print help messages:

```
iriscli upgrade query-signals --help
```

## Example

Query the current information of signals.

```
iriscli upgrade query-signals
```

```
signalsVotingPower/totalVotingPower = 0.5000000000
```

```
iriscli upgrade query-signals --detail
```

```
iva15cv33a67cfey5eze7238hck6yngw3694ak2elm   100.0000000000
siganalsVotingPower/totalVotingPower = 0.5000000000
```
