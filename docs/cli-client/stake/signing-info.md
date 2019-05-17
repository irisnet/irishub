# iriscli stake signing-info

## Description

Query a validator's signing information

## Usage

```
iriscli stake signing-info <validator-pubkey> <flags>
```

Print help messages:
```
iriscli stake signing-info --help
```

## Examples

### Query specified validator's signing information

```
iriscli stake signing-info <validator-pubkey>
```

After that, you will get specified validator's signing information.

```txt
  Signing Info
  Start Height:          0
  Index Offset:          3506
  Jailed Until:          1970-01-01 00:00:00 +0000 UTC
  Missed Blocks Counter: 0
```
