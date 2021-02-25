---
order: 3
---

# KMS - 密钥管理系统

## 简介

[Tendermint KMS](https://github.com/iqlusioninc/tmkms)是一个密钥管理服务，可将密钥管理与Tendermint节点分开。此外它还具有其他优势，例如：

- 提高安全和风险管理策略
- 统一的API和对各种HSM（硬件安全模块）的支持
- 双签保护（基于软件或硬件）

## 构建

可以在[这里](https://github.com/iqlusioninc/tmkms#installation)找到详细的构建说明

::: tip
在编译KMS的时候，确保已启用`yubihsm`功能：
:::

```bash
cargo install tmkms --features=yubihsm --version=0.10.0-beta2
```

## 初始化

为IRISHub初始化配置文件

```bash
tmkms init -n irishub /path/to/kms/home
```

## 配置

如果要启用KMS，首先需要在`<iris-home>/config/config.toml`文件中编辑`priv_validator_laddr`。例如：

```toml
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "localhost:26658"
```

然后，你需要下载[priv_validator_state.json示例](https://github.com/irisnet/irishub/blob/master/docs/tools/priv_validator_state.json)并按照`<iris-home>/data/priv_validator_state.json`修改所有字段值。

接下来，你只要只需要按如下方式编辑配置文件`/path/to/kms/home/tmkms.toml`：

- 配置`state_file`为上一步完成的`priv_validator_state.json`。
- 将你的Yubihsm密码写入文件`yubihsm-password.txt`并将`password_file`设置为它。
- 编辑`addr`指向你的`iris`实例（注意：无需指定连接ID，仅让它保持格式如：tcp://localhost:26658）。
- 调整`chain_id` 以匹配你的`<iris-home>/config/genesis.json`中的设置。
- 编辑`auth`以授权访问你的Yubihsm。
- 编辑`keys`确定你将使用的pubkey。
- 编辑`protocol_version`为v0.34。

然后启动 tmkms

```bash
tmkms start -c /path/to/kms/home/tmkms.toml
```

### 使用YubiHSM

有关如何使用YubiHSM2设置KMS的更多信息，请参阅[这里](https://github.com/iqlusioninc/tmkms/blob/master/README.yubihsm.md)。

如果要导入已存在的IRIShub private_key，可以：

```bash
tmkms yubihsm keys import <iris_home>/config/priv_validator.json -i <id> -t json -c /path/to/kms/home/tmkms.toml
```
