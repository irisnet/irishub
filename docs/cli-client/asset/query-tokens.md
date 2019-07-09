# iriscli asset query-tokens

## Description

Query the collection of tokens issued on IRIS Hub based on criteria.

## Usage

```bash
iriscli asset query-tokens [flags]
```

## Unique Flags

| Name      | Type   | Required | Default | Description                                                  |
| --------- | ------ | -------- | ------- | ------------------------------------------------------------ |
| --source  | string | false    | all     | Token Source: native / gateway / external                    |
| --gateway | string | false    |         | The unique moniker of the gateway, required when source is gateway |
| --owner   | string | false    |         | The owner of the tokens                                      |

## Query rules

- when source is native
  - gateway will be ignored
  - owner optional
- When source is gateway
  - gateway required
  - owner will be ignored (because gateway tokens are all owned by the gateway)
- when source is external
  - gateway and owner are ignored
- when the gateway is not empty
  - source optional

## Examples

### Query all tokens

```bash
iriscli asset query-tokens
```

### Query all native tokens

```bash
iriscli asset query-tokens --source=native
```

### Query all tokens of the gateway named "cats"

```bash
iriscli asset query-tokens --gateway=cats
```

### Query all tokens of the specified owner

```bash
iriscli asset query-tokens --owner=<address>
```
