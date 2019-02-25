# KMS - Key Management System

....

## What is a KMS?

...

## 构建

可以在[这里](https://github.com/irisnet/kms#installation)找到详细的构建说明。

::: 提示
在编译KMS的时候, 确保您已启用适用的功能：
:::

| Backend               | Recommended Command line              |
|-----------------------|---------------------------------------|
| YubiHSM               | ```cargo build --features yubihsm```  |
| Ledger+Tendermint App | ```cargo build --features ledgertm``` |

## Configuration

KMS提供了不同的选择

- [Using a CPU-based signer](kms_cpu.md)
- [Using a YubiHSM](kms_ledger.md)
- [Using a Ledger device running the Tendermint Validator app](kms_ledger.md)
