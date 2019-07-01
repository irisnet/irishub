# iriscli asset query-token

## Description

Query a token issued on IRIS Hub.

## Usage

```bash
iriscli asset query-token <token-id>
```

### Global Unique Token ID Generation Rule
    
- When Source is native: ID = [Symbol], e.g. iris
    
- When Source is external: ID = x.[Symbol], e.g. x.btc
    
- When Source is gateway: ID = [Gateway].[Symbol], e.g. cats.kitty

## Examples

### Query the native token named "kitty"

```bash
iriscli asset query-token kitty
```

### Query the token of gateway "cats" named "kitty"

```bash
iriscli asset query-token cats.kitty
```

### Query the external token named "btc"

```bash
iriscli asset query-token x.btc
```