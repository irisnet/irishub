# KMS - 密钥管理系统

....

## What is a KMS?

...

## 构建

可以在[这里](https://github.com/irisnet/kms#installation)找到详细的构建说明。

::: tip
在编译KMS的时候, 确保您已启用适用的功能：
:::

| Backend               | Recommended Command line              |
|-----------------------|---------------------------------------|
| YubiHSM               | ```cargo build --features yubihsm```  |
| Ledger+Tendermint App | ```cargo build --features ledgertm``` |

## Configuration

如果要启用KMS，首先需要在`~/.iris/config/config.toml`文件中编辑`priv_validator_laddr`。例如：
```text
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
Priv_validator_laddr = "localhost:26658"
```

KMS提供了不同的选择

- [Using a CPU-based signer](kms_cpu.md)
- [Using a YubiHSM](kms_yubihsm.md)
- [Using a Ledger device running the Tendermint Validator app](kms_ledger.md)
