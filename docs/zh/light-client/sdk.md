# IRISnet SDKs

IRIShub-Chain-SDK是根据IRIShub提供的API制作的一个简单的软件开发工具包，为用户快速开发基于irishub链的应用程序提供了极大的便利。

## 设计目标和概念

包客户端是整个SDK功能的入口。 SDKConfig用于配置SDK参数。

该SDK主要提供以下模块的功能，包括：auth，bank，gov，htlc，keys，nft，oracle，random，record，service，staking，token。

ClientConfig组件主要包含SDK中使用的参数，具体含义如下表所示：

| 参数     | 类型          | 描述                                             |
| --------- | ------------- | ------------------------------------------------------- |
| NodeURI   | string        | 连接到SDK的irishub节点的RPC地址，例如：localhost：26657 |
| GRPCAddr  | string        | 连接到SDK的irishub节点的GRPC地址，例如：localhost：9090 |
| Network   | enum          | irishub网络类型，值：“ Testnet”，“ Mainnet”             |
| ChainID   | string        | irishub的ChainID，例如：`irishub`                       |
| Gas       | uint64        | 交易所需支付的最大汽油费用，例如：“ 20000”              |
| Fee       | DecCoins      | 交易须支付的交易费                                      |
| KeyDAO    | KeyDAO        | 私钥管理界面，如果用户不提供，则使用默认的“ LevelDB”    |
| Mode      | enum          | 交易广播模式，值：`Sync`,`Async`, `Commit`              |
| StoreType | enum          | 私钥存储方法，值：`Keystore`,`PrivKey`                  |
| Timeout   | time.Duration | 事务超时时间，例如：5s                                  |
| Level     | string        | 日志等级例如：`info`                                    |

## 构造，签名和广播交易

如果要使用“ SDK”发送转账交易，使用“ irishub-sdk-go”的示例如下：

还有更多查询和发送交易的示例：

```go
coins, err := types.ParseDecCoins("0.1iris")
to := "iaa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
baseTx := types.BaseTx{
    From:     "username",
    Gas:      20000,
    Memo:     "test",
    Mode:     types.Commit,
    Password: "password",
}

result, err := client.Bank.Send(to, coins, baseTx)
```

查询最新区块信息

```go
block, err := client.BaseClient.Block(context.Background(),nil)
```

根据指定 TxHash查询交易:

```go
txHash := "D9280C9217B5626107DF9BC97A44C42357537806343175F869F0D8A5A0D94ADD"
txResult, err := client.BaseClient.QueryTx(txHash)
```

**注意**：如果您使用相关的API发送交易，则应实现`KeyDAO`接口。 使用NewKeyDaoWithAES方法初始化KeyDAO实例，默认情况下将使用AES加密方法。

## 私钥管理

以 irishub-sdk-go为例，接口定义如下:

```go
type KeyDAO interface {
    AccountAccess
    Crypto
}

type AccountAccess interface {
    Write(name string, store Store) error
    Read(name string) (Store,error)
    Delete(name string) error
}
type Crypto interface {
    Encrypt(data string, password string) (string, error)
    Decrypt(data string, password string) (string, error)
}
```

其中，`Store`包括两种存储方法，一种基于私钥，定义如下： 

```go
type KeyInfo struct {
    PrivKey string `json:"priv_key"`
    Address string `json:"address"`
}
```

另一种基于keystore

```go
type KeystoreInfo struct {
    Keystore string `json:"keystore"`
}
```

您可以灵活选择任何私钥管理方法。 “加密”和“解密”接口用于加密和解密密钥。 如果用户未实现，则默认为使用AES。 示例如下：

`KeyDao`实现`AccountAccess`接口：

```go
// Use memory as storage, use with caution in build environment
type MemoryDB struct {
    store map[string]Store
    AES
}

func NewMemoryDB() MemoryDB {
    return MemoryDB{
        store: make(map[string]Store),
    }
}
func (m MemoryDB) Write(name string, store Store) error {
    m.store[name] = store
    return nil
}

func (m MemoryDB) Read(name string) (Store, error) {
    return m.store[name], nil
}

func (m MemoryDB) Delete(name string) error {
    delete(m.store, name)
    return nil
}

func (m MemoryDB) Has(name string) bool {
    _, ok := m.store[name]
    return ok
}
```

## Go, JS, Java SDK 文档

IRISnet SDK的文档如下：

- [Go SDK docs](https://github.com/irisnet/irishub-sdk-go/blob/master/README.md)
- [JavaScript SDK docs](sdk-js.irisnet.org)
- Java SDK docs (敬请期待)
