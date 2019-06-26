# iriscli asset edit-token

## 描述

编辑指定ID的资产信息

## 使用方式

```shell
iriscli asset edit-token <token-id> --name=<name> --symbol-at-source=<symbol-at-source> --symbol-min-alias=<min-alias> --max-supply=<max-supply> --mintable=<mintable> --from=<your account name> --chain-id=<chain-id> --fee=0.6iris
```

## 特有的标志

| Name               | Type   | Required | Default | Description                                                  |
| ------------------ | ------ | -------- | ------- | ------------------------------------------------------------ |
| --name             | string | 否       | ""      | 资产名称，例如：IRIS Network                                 |
| --symbol-at-source | string | 否       | ""      | Source为 external 或 gateway 的时候，可以用来指定在源链上的Symbol |
| --symbol-min-alias | string | 否       | ""      | 资产最小单位别名                                             |
| --max-supply       | uint   | 否       | 0       | 以基准单位计的最大发行量                                     |
| --mintable         | bool   | 否       | false   | 初始发行后是否允许增发                                       |

## 示例

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

