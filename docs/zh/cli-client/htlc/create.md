# iriscli htlc create

## 描述

创建一个 HTLC

## 使用方式

```bash
iriscli htlc create --receiver=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --hash-lock=<hash-lock> --time-lock=<time-lock> --timestamp=<timestamp>
```

## 标志

| 命令，速记                  | 类型     | 是否必须 | 默认值 | 描述                                         |
| ------------------------- | -------- | ------ | ----- | -------------------------------------------- |
| --receiver                | string   | 是     |       | 收款人地址                                     |
| --receiver-on-other-chain | bytesHex | 是     |       | 另一条链上的收款地址                            |
| --amount                  | string   | 是     |       | 要发送的金额                                   |
| --hash-lock               | bytesHex | 是     |       | 由 Secret 和 时间戳（如果提供）生成的 sha256 哈希 |
| --time-lock               | string   | 是     |       | 资金锁定的区块数                                |
| --timestamp               | uint     |        |       | 参与生成 hashlock 的10位时间戳（可选）           |

## 示例

### 创建一个 HTLC

```bash
iriscli htlc create \
--from=userX \
--receiver=faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5 \
--receiver-on-other-chain=bb9188215a6112a6f7eb93e3e929197b3d44004cb691f95babde84cc18789364 \
--hash-lock=e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561 \
--amount=10iris \
--time-lock=50 \
--timestamp=1580000000 \
--fee=0.3iris \
--chain-id=testNet \
--commit
```
