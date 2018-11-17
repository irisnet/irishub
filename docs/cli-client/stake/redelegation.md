# iriscli stake redelegation

## Description

Query a redelegation record based on delegator and a source and destination validator address

## Usage

```
iriscli stake redelegation [flags]
```

## Flags

| Name, shorthand            | Default                    | Description                                                         | Required |
| -------------------------- | -------------------------- | ------------------------------------------------------------------- | -------- | 
| --address-delegator        |                            | [string] Bech address of the delegator                              | Yes      |
| --address-validator-dest   |                            | [string] Bech address of the destination validator                  | Yes      |
| --address-validator-source |                            | [string] Bech address of the source validator                       | Yes      |
| --chain-id                 |                            | [string] Chain ID of tendermint node                                |          |
| --height                   | most recent provable block | block height to query                                               |          |
| --help, -h                 |                            | help for redelegation                                               |          |
| --indent                   |                            | Add indent to JSON response                                         |          |
| --ledger                   |                            | Use a connected Ledger device                                       |          |
| --node                     | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --trust-node               | true                       | Don't verify proofs for responses                                   |          |

## Examples

### Query a redelegation record

```shell
iriscli stake redelegation --address-validator-source=SourceValidatorAddress --address-validator-dest=DestinationValidatorAddress --address-delegator=DelegatorAddress
```

After that, you will get specified redelegation's info based on delegator and a source and destination validator address

```txt
Redelegation
Delegator: faa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk
Source Validator: fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2
Destination Validator: fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd
Creation height: 1130
Min time to unbond (unix): 2018-11-16 07:22:48.740311064 +0000 UTC
Source shares: 0.1000000000
Destination shares: 0.1000000000
```
