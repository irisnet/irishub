# iriscli record download

## Description

Download related data with unique record ID to specified file

## Usage

```
iriscli record download [record ID] [flags]
```

## Flags

| Name, shorthand | Default                    | Description                                                       | Required |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                              | Yes      |
| --file-name     |                            | [string] Download file name                                       | Yes      |
| --height        | most recent provable block | [int] block height to query                                       |          |
| --help, -h      |                            | help for download                                                 |          |
| --indent        |                            | Add indent to JSON response                                       |          |
| --ledger        |                            | Use a connected Ledger device                                     |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain |          |
| --record-id     |                            | [string] record ID                                                |          |
| --trust-node    | true                       | Don't verify proofs for responses                                 |          |

## Examples

### Download a record

```shell
iriscli record download --chain-id=test --record-id=MyRecordID --file-name="download.txt"
```

After that, you will get download file under iriscli home directory with the specfied record ID info.

```txt
[ONCHAIN] Downloading ~/.iriscli/download.txt from blockchain directly...
[ONCHAIN] Download file from blockchain complete.
```
