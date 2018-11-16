# Record User Guide

## Basic Function Description

1. On-Chain record for the text data
2. On-Chain record for the text file (TODO)
3. On-Chain governance for the record params (TODO)

## Interaction Process

### Record process

1. Any users can initiate a record request. It will cost you some tokens. If there’s no record of the data on the existing chains, the request will be completed successfully and the relevant metadata will be recorded onchain. And you will be returned a record ID to confirm your ownership of the data.
2. If any others initiate a record request for the same data, the request will be directly rejected and it will hint that the relevant record data has already existed.
3. Any users can search/download onchain based on the record ID.
4. At present, the maximum amount of stored data at most 1K Bytes. In the future, the dynamic adjustment of parameters will be implemented in conjunction with the governance module.

## Usage Scenarios

### Build Usage Scenarios

```
rm -rf iris
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=record-test -o --home=iris
iris start --home=iris
```

### Usage Scenarios of Record on Chains

Scenario 1: Record the data through cli

```
# Specify the text data which needs to be recorded by --onchain-data

iriscli record submit --description="test" --onchain-data=x --from=x --fee=0.04iris

# Result
Committed at block 4 (tx hash: F649D5465A28842B50CAE1EE5950890E33379C45, response: {Code:0 Data:[114 101 99 111 114 100 58 97 98 57 100 57 57 100 48 99 102 102 54 53 51 100 99 54 101 56 53 52 53 99 56 99 99 55 50 101 53 53 51 51 100 101 97 97 97 49 50 53 53 50 53 52 97 100 102 100 54 98 48 48 55 52 101 50 56 54 57 54 54 49 98] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3857 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 114 101 99 111 114 100]} {Key:[111 119 110 101 114 65 100 100 114 101 115 115] Value:[102 97 97 49 109 57 51 99 103 57 54 51 56 121 115 104 116 116 100 109 119 54 57 97 121 118 51 103 106 53 102 116 116 109 108 120 51 99 102 121 107 109]} {Key:[114 101 99 111 114 100 45 105 100] Value:[114 101 99 111 114 100 58 97 98 57 100 57 57 100 48 99 102 102 54 53 51 100 99 54 101 56 53 52 53 99 56 99 99 55 50 101 53 53 51 51 100 101 97 97 97 49 50 53 53 50 53 52 97 100 102 100 54 98 48 48 55 52 101 50 56 54 57 54 54 49 98]} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[2 189 149 142 250 208 0]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-record",
     "completeConsumedTxFee-iris-atto": "\u0002\ufffd\ufffd\ufffd\ufffd\ufffd\u0000",
     "ownerAddress": "faa1m93cg9638yshttdmw69ayv3gj5fttmlx3cfykm",
     "record-id": "record:ab9d99d0cff653dc6e8545c8cc72e5533deaaa1255254adfd6b0074e2869661b"
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
