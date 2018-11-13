# iriscli record query

## Description

Query details of the existing record

## Usage

```
iriscli record query [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                       | Required |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                              | Yes      |
| --height        | most recent provable block | block height to query                                             |          |
| --help, -h      |                            | help for query                                                    |          |
| --indent        |                            | Add indent to JSON response                                       |          |
| --ledger        |                            | Use a connected Ledger device                                     |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --record-id     |                            | [string] Record ID for query                                      | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                 |          |

## Examples

### Query a record

```shell
iriscli record query --chain-id=test --record-id=MyRecordID
```

After that, you will get detail info for the record which has the specfied record ID.

```txt
{
  "submit_time": "2018-11-13 15:31:36",
  "owner_addr": "faa122uzzpugtrzs09nf3uh8xfjaza59xvf9rvthdl",
  "record_id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6",
  "description": "description",
  "data_hash": "ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6",
  "data_size": "24",
  "data": "this is my on chain data"
}
```
