# Cosmos SDK v0.53.8 升级

本分支使用 Go 模块路径 `github.com/irisnet/irishub/v5`，升级到 Cosmos SDK
`v0.53.8` 和 IBC-Go `v10.7.0`。最低 Go 版本为 `1.23.8`。

## 固定依赖

依赖通过伪版本固定到以下分支提交，构建时不依赖可变的分支头：

| 依赖 | 分支 | 提交 |
| --- | --- | --- |
| nft-transfer | `feat/cosmos-sdk-v0.53.8` | `87faf8d2b44a` |
| tibc-go | `feat/cosmos-sdk-v0.53.8` | `0b8e7b4af967` |
| mods.irisnet.org/modules/* | `feat/cosmos-sdk-v0.53.8` | `f4182d6e7472` |
| bianjieai/ethermint | `release/v0.24.0` | `4a4e53b81b98` |

IBC-Go 直接使用 `v10.7.0`，不通过 `replace` 降级。按照 Ethermint 分支的替换配置，
实际 EVM 实现保留在 go-ethereum `v1.10.26`。可用
`go list -m -json github.com/cosmos/ibc-go/v10 github.com/ethereum/go-ethereum`
检查实际版本、替换信息和源码目录。

## 链上升级

本次注册的软件升级计划名称是 **`cosmos-sdk-v0.53.8`**。现有链需要在约定的升级
高度执行该计划；仅替换二进制并直接打开旧数据库不能代替状态迁移。

升级加载器删除不再使用的 `capability` KV store；升级处理器移除其模块版本记录，
再运行已注册的 SDK 和 IBC 模块迁移，包括 IBC transfer 的 DenomTrace 到 Denom
格式迁移。原有 transfer、NFT transfer 和 ICA host 路由继续注册；Tendermint
和 solo-machine 客户端按 IBC v10 要求显式注册。IBC 客户端恢复、参数修改等
治理操作使用 v10 消息，由 `gov` 模块作为 authority 执行。

`auth` 的 PreBlocker 排在 `upgrade` 之后。未启用无序交易、epochs 或 protocolpool。

本次不保留 LSM 功能，已移除旧 v3 升级处理器中启用 LSM 的步骤。普通 staking 的
protobuf 字段编号保持兼容，LSM 附加字段不再由上游类型读取；原 LSM 使用的
`0x81` 至 `0x87` 存储前缀在 SDK v0.53.8 中明确保留，不会复用。本次没有销毁、
赎回或重新分配任何存量 LSM 代币，也没有改写已有验证人的最低自委托值。

## Docker 构建

Dockerfile 默认使用 `golang:1.24.9-alpine` 和 `alpine:3.18`。Docker Hub 不可用时，
可通过构建参数切换到镜像缓存：

```shell
docker build -t irishub:sdk0538 --build-arg EVM_CHAIN_ID=6688 --build-arg BUILDER_IMAGE=mirror.gcr.io/library/golang:1.24.9-alpine --build-arg RUNTIME_IMAGE=mirror.gcr.io/library/alpine:3.18 .
```

## 验证范围

新增测试覆盖真实应用创世初始化、首块执行、IBC 路由，以及通过已注册升级处理器
迁移旧 DenomTrace 编码并保持 IBC 资产 denom 不变。仓库现有测试用于检查本地
模块和 CLI 的兼容性。

这些检查没有使用生产链数据库快照。生产升级前仍需用目标链快照演练升级高度、
已有 IBC 通道/客户端和存量 LSM 资产的处置；普通单元测试不能证明这些链上状态
已完成验证。

参考：[SDK 升级说明](https://github.com/cosmos/cosmos-sdk/blob/v0.53.8/UPGRADING.md)、
[IBC v8.1 到 v10 升级说明](https://github.com/cosmos/ibc-go/blob/v10.7.0/docs/docs/05-migrations/13-v8_1-to-v10.md)。
