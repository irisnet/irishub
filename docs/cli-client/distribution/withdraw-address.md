# iriscli distribution withdraw-address

## Description

Query the withdraw address of a delegator

## Usage

```
iriscli distribution withdraw-address <delegator-address> <flags>
```

Print help messages:

```
iriscli distribution withdraw-address --help
```

## Examples

```
iriscli distribution withdraw-address <delegator-address>
```

Example response:
```text
iaa1ezzh0humhy3329xg4avhcjtay985nll06lgq50
```
If the given delegator didn't specify other withdraw address, the query result will be empty.