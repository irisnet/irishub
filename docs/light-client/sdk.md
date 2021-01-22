# IRISnet SDKs

IRIShub Chain SDK makes a simple package of API provided by IRIShub, which provides great convenience for users to quickly develop applications based on irishub chain.

## Design Goals and Concepts

Package client is the entrance of the entire SDK function. SDKConfig is used to configure SDK parameters.

The SDK mainly provides the functions of the following modules, including: auth, bank, gov, htlc, keys, nft, oracle, random, record, service, staking, token.

The `ClientConfig` component mainly contains the parameters used in the SDK, the specific meaning is shown in the table below:

| Iterm     | Type          | Description                                                                                           |
| --------- | ------------- | ----------------------------------------------------------------------------------------------------- |
| NodeURI   | string        | The RPC address of the irishub node connected to the SDK, for example: localhost: 26657               |
| GRPCAddr  | string        | The GRPC address of the irishub node connected to the SDK, for example: localhost: 9090               |
| Network   | enum          | irishub network type, value: `Testnet`,`Mainnet`                                                      |
| ChainID   | string        | ChainID of irishub, for example: `irishub`                                                            |
| Gas       | uint64        | The maximum gas to be paid for the transaction, for example: `20000`                                  |
| Fee       | DecCoins      | Transaction fees to be paid for transactions                                                          |
| KeyDAO    | KeyDAO        | Private key management interface, If the user does not provide it, the default `LevelDB` will be used |
| Mode      | enum          | Transaction broadcast mode, value: `Sync`,`Async`, `Commit`                                           |
| StoreType | enum          | Private key storage method, value: `Keystore`,`PrivKey`                                               |
| Timeout   | time.Duration | Transaction timeout, for example: `5s`                                                                |
| Level     | string        | Log output level, for example: `info`                                                                 |

## Generating, Signing and Broadcasting Transactions

If you want to use `SDK` to send a transfer transaction, the example with `irishub-sdk-go` is as follows:

There is more example of query and send tx:

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

Query Latest Block info:

```go
block, err := client.BaseClient.Block(context.Background(),nil)
```

Query Tx from specify TxHash:

```go
txHash := "D9280C9217B5626107DF9BC97A44C42357537806343175F869F0D8A5A0D94ADD"
txResult, err := client.BaseClient.QueryTx(txHash)
```

**Note**: If you use the relevant API for sending transactions, you should implement the `KeyDAO` interface. Use the `NewKeyDaoWithAES` method to initialize a `KeyDAO` instance, which will use the `AES` encryption method by default.

## Private Key Management

Take irishub-sdk-go as an example, the interface definition is as follows:

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

Among them, `Store` includes two storage methods, one is based on the private key, which is defined as follows:

```go
type KeyInfo struct {
    PrivKey string `json:"priv_key"`
    Address string `json:"address"`
}
```

The other is based on the keystore, defined as follows:

```go
type KeystoreInfo struct {
    Keystore string `json:"keystore"`
}
```

You can flexibly choose any of the private key management methods. The `Encrypt` and `Decrypt` interfaces are used to encrypt and decrypt the key. If the user does not implement it, the default is to use `AES`. Examples are as follows:

`KeyDao` implements the `AccountAccess` interface:

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

## Go, JS, Java SDK Docs

The docs of IRISnet SDKs are as follows:

- [Go SDK docs](https://github.com/irisnet/irishub-sdk-go/blob/master/README.md)
- [JavaScript SDK docs](https://github.com/irisnet/irishub-sdk-js/blob/master/README.md)
- [Java SDK docs](https://github.com/irisnet/irishub-sdk-java/blob/master/README.md)
