# iriscli asset transfer-token-owner

## Description

Transfer control of the token

## Usage

```shell
iriscli asset transfer-token-owner <token-id> --to=<new owner>
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --to           | string | Yes | "" | Destination address |

## Example

1. Generate transaction file

   ```shell
   iriscli asset transfer-token-owner neo --to=faa1t5a294zz2xkk6ulz2pz8uwppv3fexydf90smg8 --from=node0 --chain-id=irishub-test --fee=0.4iris > src.json
   ```

2. Current token owner signature

   ```shell
   iriscli tx sign src.json --name=node0 --chain-id=irishub-test > tx1.json
   ```

3. Destination owner signature

   ```shell
   iriscli tx sign tx1.json --name=neoOwner --chain-id=irishub-test > tx2.json
   ```

   You can view file(tx2.json) content,as follows:

   ```json
   {
     "type": "irishub/bank/StdTx",
     "value": {
       "msg": [
         {
           "type": "irishub/asset/MsgTransferTokenOwner",
           "value": {
             "src_owner": "faa1sf4xrfq3p45hlelp5pnw5rf4llsfg4st07mhjc",
             "dst_owner": "faa1t5a294zz2xkk6ulz2pz8uwppv3fexydf90smg8",
             "token_id": "neo"
           }
         }
       ],
       "fee": {
         "amount": [
           {
             "denom": "iris-atto",
             "amount": "400000000000000000"
           }
         ],
         "gas": "50000"
       },
       "signatures": [
         {
           "pub_key": {
             "type": "tendermint/PubKeySecp256k1",
             "value": "A2RHE/FJGsg4NMzAgoHUrUpfm+wV/unz4Jm5BF/juE68"
           },
           "signature": "wU04prBDVQLo1dT3tXtpHTSKUfCQ+zDiBZg921PL94YRNdd7pKFx/x15ltnDK5HN5Ah0feJjHkQiHXJaZ7+ncA==",
           "account_number": "2",
           "sequence": "23"
         },
         {
           "pub_key": {
             "type": "tendermint/PubKeySecp256k1",
             "value": "A2kjyJhgOJKk+AVDf9u8EnKQ7zr1vCwDwGIpmQGk325D"
           },
           "signature": "vdO8o4O5Fqv1Lkp2Cge93RZx+ODS8Dbt893gubc4NkhTxjgs+2Yt6YUBydInpfmJFwXtvjmGUPzW9Kgd8kjGig==",
           "account_number": "5",
           "sequence": "0"
         }
       ],
       "memo": ""
     }
   }
   ```

4. Broadcast transaction

   ```shell
   iriscli tx broadcast tx2.json
   ```