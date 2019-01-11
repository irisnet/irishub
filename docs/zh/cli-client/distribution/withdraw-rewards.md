# iriscli distribution withdraw-rewards

## 介绍

取回收益

## 用法

```
iriscli distribution withdraw-rewards [flags]
```

打印帮助信息:

```
iriscli distribution withdraw-rewards --help
```
## 标志位

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --chain-id       | String | 否       |                       | tendermint 节点网络ID                                        |
| --account-number | int    | 否       |                       | 账户数字用于签名通证发送                                     |
| --async          |        | 否       | True                  | 异步广播传输信息                                             |
| --dry-run        |        | 否       |                       | 忽略--gas 标志 ，执行仿真传输，但不广播。                    |
| --fee            | String | 是       |                       | 设置传输需要的手续费                                         |
| --from           | String | 是       |                       | 用于签名的私钥名称                                           |
| --gas            | String | 否       | 20000                 | 每笔交易设定的gas限额; 设置为“simulate”以自动计算所需气体    |
| --gas-adjustment | Float  | 否       | 1                     | 调整因子乘以传输模拟返回的估计值; 如果手动设置气体限制，则忽略该标志 |
| --generate-only  |        | 否       |                       | 创建一个未签名的传输并写到标准输出中。                       |
| --indent         |        | 否       |                       | 在JSON响应中增加缩进                                         |
| --json           |        | 否       |                       | 以json格式返回输出                                           |
| --memo           | String | 否       |                       | 传输中的备注信息                                             |
| --print-response |        | 否       |                       | 返回传输响应 (仅仅当 async = false时有效)                    |
| --sequence       | Int    | 否       |                       | 等待签名传输的序列号。                                       |
| --ledger         | String | 否       |                       | 使用一个联网的分账设备                                       |
| --node           | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。                    |
| --trust-node     | String | 否       | True                  | 不验证响应的证明                                             |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|

## 特有标志位

| 名称                | 类型   | 是否必填 | 默认值  | 功能描述        |
| --------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --only-from-validator | string | false    | ""       | 验证人地址，仅取回在这个验证人上的委托收益 |
| --is-validator        | bool   | false    | false    | 取回验证人佣金收益 |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|

不能同时使用两个flags。

## 示例

1. 仅取回某一个委托产生的收益
    ```
    iriscli distribution withdraw-rewards --only-from-validator fva134mhjjyyc7mehvaay0f3d4hj8qx3ee3w3eq5nq --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
2. 取回所有委托产生的收益，不包含验证人的佣金收益:
    ```
    iriscli distribution withdraw-rewards --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
3. 取回所有委托产生的收益以及验证人的佣金收益:
    ```
    iriscli distribution withdraw-rewards --is-validator=true --from mykey --fee=0.004iris --chain-id=irishub-test
    ```