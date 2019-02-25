# Setting up Tendermint KMS + YubiHSM

::: danger Warning
以下说明是一个简短的演练，而不是一个全面的指南。
:::

::: danger Warning
KMS目前正在进行中。细节可能有所不同请谨慎使用，风险自负。
:::

## YubiHSM
[YubiHSM](https://www.yubico.com/products/yubihsm/): hardware security module providing root of trust for servers and computing devices.

## KMS配置

在本节中，我们将配置KMS以使用YubiHSM。

#### 配置文件

可以在[这里](https://github.com/irisnet/kms/blob/master/tmkms.toml.example)找到其他配置示例。

- 使用以下内容创建一个`~/.tmkms/tmkms.toml`文件：

```toml
# Example KMS configuration file
[[validator]]
addr = "tcp://localhost:26658"    # or "unix:///path/to/socket"
chain_id = "fuxi"
reconnect = true # true is the default
secret_key = "~/.tmkms/secret_connection.key"
[[providers.yubihsm]]
adapter = { type = "usb" }
auth = { key = 1, password = "password" } # Default YubiHSM admin credentials. Change ASAP!
keys = [{ id = "test", key = 1 }]
#serial_number = "0123456789" # identify serial number of a specific YubiHSM to connect to
```

- 编辑 `addr` 指向你的 `iris` 实例。
- 调整 `chain-id` 以匹配你的 `.iris/config/config.toml` 设置。
- 编辑 `auth` 以授权访问你的yubihsm。
- 编辑 `keys` 确定您将使用哪个pubkey。[如何导入密钥？](#import-private-key-to-yubihsm))

#### Import private key to yubihsm

```bash
tmkms yubihsm keys import  -p ~/.iris/config/priv_validator.json [id]
```

#### 生成连接密钥

现在你需要生成secret_key

```bash
tmkms keygen ~/.tmkms/secret_connection.key
```

#### 启动KMS

最后一步是检索将在`iris`中使用的验证人密钥。

启动KMS:

```bash
tmkms start -c ~/.tmkms/tmkms.toml
```

The output should look similar to:

```text
07:28:24 [INFO] tmkms 0.3.0 starting up...
07:28:24 [INFO] [keyring:ledgertm:ledgertm] added validator key icp1zcjduepqa9y67dqgug4u4stf5sf0arvjrnty8eenlfj22vnh78cmejd8qdss8t6ljg
07:28:24 [INFO] KMS node ID: 1BC12314E2E1C29015B66017A397F170C6ECDE4A
```