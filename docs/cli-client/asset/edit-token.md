# iriscli asset edit-token

## Description

edit token's information by token-id

## Usage

```shell
iriscli asset edit-token <token-id> --name=<name> --symbol-at-source=<symbol-at-source> --symbol-min-alias=<min-alias> --max-supply=<max-supply> --mintable=<mintable> --from=<your account name> --chain-id=<chain-id> --fee=0.6iris
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --name           | string | No | "" | the asset name, e.g. IRIS Network |
| --symbol-at-source | string | No | "" | the source symbol of a gateway or external asset |
| --symbol-min-alias | string | No | "" | the asset symbol minimum alias |
| --max-supply | uint | No | 0 | the max supply token of asset |
| --mintable | bool | No | false | whether the asset can be minted, default false |


## Example

```shell
iriscli asset edit-token eth --name="ETH TOKEN" --symbol-at-source="ETH" --symbol-min-alias=atto --max-supply=100000000000 --mintable=true --from=node0 --chain-id=irishub-test --fee=0.4iris  --home=$iris_root_path --commit
```

输出信息:
```txt
Password to sign with 'node0':
Committed at block 502 (tx hash: 3D131F2D1E0B200206E8023E70C9442142DA27EBC42675451E39BF84B6343C6F, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 5350,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "edit_token"
     },
     {
       "key": "token-id",
       "value": "eth"
     }
   ]
 })
```
