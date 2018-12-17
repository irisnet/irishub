# iriscli gov deposit

## Description
 
Deposit tokens for activing proposal
 
## Usage
 
```
iriscli gov deposit [flags]
```

Print help messages:

```
iriscli gov deposit --help
```
## Flags
 
| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] Deposit of proposal                                                                                                                         | Yes      |
| --proposal-id    |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
## Examples

### Deposit

```shell
iriscli gov deposit --chain-id=test --proposal-id=1 --deposit=50iris --from=node0 --fee=0.01iris
```

After you enter the correct password, you could deposit 50iris to make your proposal active which can be voted, after you enter the correct password, you're done with depositing iris tokens for an activing proposal.

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
     "depositor": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "proposal-id": "1",
     "voting-period-start": "1"
   }
 })
```
