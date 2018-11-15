# Record User Guide

## Basic Function Description

1. On-Chain record for the text data
2. On-Chain record for the text file (TODO)
3. On-Chain governance for the record params (TODO)

## Interaction

### Record Introduction

Record metadata components

| Fields in record metadata    | Description                              |
| ---------------------------- | ---------------------------------------- |
| record submit time           | data submit time                         |
| record owner address         | data owner's address on the target chain |
| record ID                    | record index ID                          |
| record description           | data description                         |
| data hash                    | uploaded data hash                       |
| data size                    | uploaded data size                       |
| data                         | uploaded data itself                     |

1. Record module generates related metadata from uploaded data on the target chain, and it has the obvious advantages such as efficiency, auditability and traceability over the traditional record technologies.
2. Any users can initiate a record request. It will cost you some tokens. If there’s no identical data on the targt chain, the request will be completed successfully and the metadata of record payload will be recorded on the target chain. And you will be returned a record ID to confirm your ownership of the data.
3. If any others initiate a record request for the same data, the request will be directly rejected and it will hint that the relevant record data has already existed.
4. Any users can search/download onchain based on the record ID.
5. At present, the maximum amount of stored data at most 1K Bytes. In the future, the dynamic adjustment of parameters will be implemented in conjunction with the governance module.

## Usage Scenarios

### Usage Scenarios of Record on Chains

Scenario 1: Record the data through cli

```
# Specify the text data which needs to be recorded by --onchain-data

iriscli record submit --description="test" --onchain-data=x --from=x --fee=0.04iris

# Result
Committed at block 1845 (tx hash: E620F3CD62BD9128443BA168296FFECC9BE2AF8F45CF21FD8FDA609DEFA253ED, response: {Code:0 Data:[114 101 99 111 114 100 58 101 56 51 102 57 102 50 55 55 100 101 99 102 54 98 57 52 102 100 55 50 55 100 52 54 53 99 49 51 52 54 57 102 53 53 50 48 51 57 56 56 57 102 98 101 102 99 56 97 49 99 55 52 101 48 54 99 54 56 48 57 98 48 101] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3978 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 114 101 99 111 114 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[111 119 110 101 114 65 100 100 114 101 115 115] Value:[102 97 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 48 117 115 121 50 50] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 111 114 100 45 105 100] Value:[114 101 99 111 114 100 58 101 56 51 102 57 102 50 55 55 100 101 99 102 54 98 57 52 102 100 55 50 55 100 52 54 53 99 49 51 52 54 57 102 53 53 50 48 51 57 56 56 57 102 98 101 102 99 56 97 49 99 55 52 101 48 54 99 54 56 48 57 98 48 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 55 57 53 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-record",
     "completeConsumedTxFee-iris-atto": "\"795600000000000\"",
     "ownerAddress": "faa15grv3xg3ekxh9xrf79zd0w077krgv5xf0usy22",
     "record-id": "record:e83f9f277decf6b94fd727d465c13469f552039889fbefc8a1c74e06c6809b0e"
   }
}

# Query the Status of the Records
iriscli record query --record-id=x

# Download the Record
iriscli record download --record-id=x --file-name="download"

```

Scenario 2: Query the transactions including recorded data onchain through cli

```
# Query the status of the records onchain
iriscli tendermint txs --tag "action='submit-record'"
```

## Details of cli

```
iriscli record submit --description="test" --onchain-data=x --chain-id="record-test" --from=x --fee=0.04iris
```

* `--onchain-data`  The data needs to be recorded


```
iriscli record query --record-id=x --chain-id="record-test"
```

* `--record-id` Record ID to be queried


```
iriscli record download --record-id=x --file-name="download" --chain-id="record-test"
```

* `--file-name` The filename of recorded data, in the directory specified by `--home`
