---
order: 3
---

# KMS - 密钥管理系统

## 简介

KMS是密钥管理系统的简称，请参阅[Tendermint KMS](https://github.com/tendermint/kms)获取更多详细信息。

## 构建

可以在[这里](https://github.com/tendermint/kms#installation)找到详细的构建说明。

::: tip
在编译KMS的时候, 确保您已启用适用的功能：
:::

| Backend               | Recommended Command line              |
|-----------------------|---------------------------------------|
| YubiHSM               | ```cargo build --features yubihsm```  |
| Ledger+Tendermint App | ```cargo build --features ledgertm``` |

## 配置

[Tendermint KMS](https://github.com/tendermint/kms)支持所有基于[tendermint](https://github.com/tendermint/tendermint)共识引擎构建的区块链，包括IRIShub。

如果要启用KMS，首先需要在`<iris-home>/config/config.toml`文件中编辑`priv_validator_laddr`。例如：

```toml
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "localhost:26658"
```

你可以下载[示例配置文件](https://github.com/tendermint/kms/blob/master/tmkms.toml.example)，其中包含对IRIShub的支持，你只需要做如下修改：

- 编辑`addr`指向你的`iris`实例。
- 调整`chain-id` 以匹配你的`<iris-home>/config/genesis.json`中的设置。
- 编辑`auth`以授权访问你的yubihsm。
- 编辑`keys`确定你将使用的pubkey。

然后启动tmkms:

```bash
tmkms start
```

KMS提供了多种选择

### 使用YubiHSM

有关如何使用YubiHSM2设置KMS的更多信息，请参阅[这里](https://github.com/tendermint/kms/blob/master/README.yubihsm.md)。

如果要导入已存在的IRIShub private_key，可以：

```bash
tmkms yubihsm keys import <iris_home>/config/priv_validator.json -i <id>
```

### 使用运行Tendermint app的ledger设备

- [Using a Ledger device running the Tendermint Validator app](kms/kms_ledger.md)
