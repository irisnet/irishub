# iriscli gov deposit

## Description
 
Deposit tokens for active proposal
 
## Usage
 
```
iriscli gov deposit <flags>
```

Print help messages:

```
iriscli gov deposit --help
```
## Flags
 
| Name, shorthand  | Default                    | Description                                                                                                                                         | Required |
| ---------------- | -------------------------- | --------------------------------------------------------------------------------------------------------------- | ---------- |
| --deposit        |                                       | Deposit of proposal                                                                                                                 | Yes          |
| --proposal-id  |                                       | ProposalID of proposal depositing on                                                                                    | Yes      |

## Examples

### Deposit

```shell
iriscli gov deposit --chain-id=<chain-id> --proposal-id=<proposal-id> --deposit=50iris --from=<key_name> --fee=0.3iris
```

```txt
Committed at block 7 (tx hash: C1156A7D383492AE5C2EB1BADE0080C3A36BE8AED491DC5B2331056BED5D60DC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7944,
   "codespace": "",
   "tags": {
     "action": "deposit",
     "depositor": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "proposal-id": "1",
     "voting-period-start": "1"
   }
 })
```

When the total deposit amount exceeds `MinDeposit`, the proposal enter the voting procedure. 

| GovParams | Critical | Important | Normal |
| ------ | ------ | ------ | ------|
| MinDeposit | 4000 iris | 2000 iris | 1000 iris |


### How to query deposit

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)
