# gRPC Client

IRIShub v1.0.0（依赖Cosmos-SDK v0.41）引入了 Protobuf 作为主要的[编码](https://github.com/cosmos/cosmos-sdk/blob/master/docs/core/encoding.md)库，这带来了可插入 SDK 的各种基于 Protobuf 的工具。一种这样的工具是 [gRPC](https://grpc.io)，这是一种现代的开源高性能 RPC 框架，具有多语言客户端支持。

## gRPC 服务端口、激活方式和配置

`grpc.Server` 是一个具体的 gRPC 服务，它产生并服务任何gRPC请求。可以在 `~/.iris/config/app.toml` 中配置：

- `grpc.enable = true|false` 字段定义了 gRPC 服务是否可用，默认为 `true`。
- `grpc.address = {string}` 字段定义了服务绑定的地址（实际上是端口，因为主机必须保持为 `0.0.0.0`），默认为 `0.0.0.0:9000`。

gRPC 服务启动后，您可以使用 gRPC 客户端向其发送请求。

## gRPC 端点

IRIShub 附带的所有可用 gRPC 端点的概述见[Protobuf 文档](./proto-docs.md)。

## 构造、签名和广播交易

可以使用 Cosmos SDK 的 `TxBuilder` 接口，通过 Golang 以编程方式处理交易。

### 构造一个交易

在生成交易之前，需要创建一个新的 `TxBuilder` 实例。 由于 SDK 支持 Amino 和 Protobuf 交易，因此第一步将是确定要使用哪种编码方案。无论您使用的是 Amino 还是 Protobuf，所有后续步骤均保持不变，因为 `TxBuilder` 抽象了编码机制。在以下代码段中，我们将使用 Protobuf。

```go
import (
    "github.com/cosmos/cosmos-sdk/simapp"
)

func sendTx() error {
    // Choose your codec: Amino or Protobuf. Here, we use Protobuf, given by the following function.
    encCfg := simapp.MakeTestEncodingConfig()

    // Create a new TxBuilder.
    txBuilder := encCfg.TxConfig.NewTxBuilder()

    // --snip--
}
```

我们还可以设置一些密钥和地址来发送和接收交易。在此，出于本教程的目的，我们将使用一些虚拟数据来创建密钥。

```go
import (
    "github.com/cosmos/cosmos-sdk/testutil/testdata"
)

priv1, _, addr1 := testdata.KeyTestPubAddr()
priv2, _, addr2 := testdata.KeyTestPubAddr()
priv3, _, addr3 := testdata.KeyTestPubAddr()
```

可以通过这些[方法](https://github.com/cosmos/cosmos-sdk/blob/v0.41.0/client/tx_config.go#L32-L45)来配置 `TxBuilder`：

```go
import (
    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func sendTx() error {
    // --snip--

    // Define two x/bank MsgSend messages:
    // - from addr1 to addr3,
    // - from addr2 to addr3.
    // This means that the transactions needs two signers: addr1 and addr2.
    msg1 := banktypes.NewMsgSend(addr1, addr3, types.NewCoins(types.NewInt64Coin("atom", 12)))
    msg2 := banktypes.NewMsgSend(addr2, addr3, types.NewCoins(types.NewInt64Coin("atom", 34)))

    err := txBuilder.SetMsgs(msg1, msg2)
    if err != nil {
        return err
    }

    txBuilder.SetGasLimit(...)
    txBuilder.SetFeeAmount(...)
    txBuilder.SetMemo(...)
    txBuilder.SetTimeoutHeight(...)
}
```

至此，以 `TxBuilder` 为基础的交易已经准备好进行签名。

### 签名一个交易

我们将编码设置为使用 Protobuf，默认情况下将使用 `SIGN_MODE_DIRECT`。 根据[ADR-020](https://github.com/cosmos/cosmos-sdk/blob/v0.41.0/docs/architecture/adr-020-protobuf-transaction-encoding.md)，每个签名者都需要对所有其他签名者的 `SignerInfo`s 进行签名。这意味着我们需要依次执行两个步骤：

- 对于每个签名者，在 `TxBuilder` 中设置签名者的 `SignerInfo`，
- 设置所有 `SignerInfo` 之后，每个签名者对 `SignDoc`（要签名的有效数据）进行签名。

在当前的 `TxBuilder` API中，两个步骤都使用相同的方法 `SetSignatures()` 完成。当前的 API 要求我们首先循环执行带不带签名的 `SetSignatures()`，仅设置 `SignerInfo`s，然后进行第二轮 `SetSignatures()` 来对正确的有效数据进行签名。

```go
import (
    cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
    "github.com/cosmos/cosmos-sdk/types/tx/signing"
    xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func sendTx() error {
    // --snip--

    privs := []cryptotypes.PrivKey{priv1, priv2}
    accNums:= []uint64{..., ...} // The accounts' account numbers
    accSeqs:= []uint64{..., ...} // The accounts' sequence numbers

    // First round: we gather all the signer infos. We use the "set empty
    // signature" hack to do that.
    var sigsV2 []signing.SignatureV2
    for i, priv := range privs {
        sigV2 := signing.SignatureV2{
            PubKey: priv.PubKey(),
            Data: &signing.SingleSignatureData{
                SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
                Signature: nil,
            },
            Sequence: accSeqs[i],
        }

        sigsV2 = append(sigsV2, sigV2)
    }
    err := txBuilder.SetSignatures(sigsV2...)
    if err != nil {
        return err
    }

    // Second round: all signer infos are set, so each signer can sign.
    sigsV2 = []signing.SignatureV2{}
    for i, priv := range privs {
        signerData := xauthsigning.SignerData{
            ChainID:       chainID,
            AccountNumber: accNums[i],
            Sequence:      accSeqs[i],
        }
        sigV2, err := tx.SignWithPrivKey(
            encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
            txBuilder, priv, encCfg.TxConfig, accSeqs[i])
        if err != nil {
            return nil, err
        }

        sigsV2 = append(sigsV2, sigV2)
    }
    err = txBuilder.SetSignatures(sigsV2...)
    if err != nil {
        return err
    }
}
```

现在已经正确配置了 `TxBuilder`。 要打印它，您可以使用初始编码配置 `encCfg` 中的 `TxConfig` 接口：

```go
func sendTx() error {
    // --snip--

    // Generated Protobuf-encoded bytes.
    txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
    if err != nil {
        return err
    }

    // Generate a JSON string.
    txJSONBytes, err := encCfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
    if err != nil {
        return err
    }
    txJSON := string(txJSONBytes)
}
```

### 广播一个交易

广播交易的首选方法是使用 gRPC，尽管也可以使用 REST（通过 `gRPC-gateway`）或 Tendermint RPC。 本教程中，我们仅介绍 gRPC 方法。

```go
import (
    "context"
    "fmt"

    "google.golang.org/grpc"

    "github.com/cosmos/cosmos-sdk/types/tx"
)

func sendTx(ctx context.Context) error {
    // --snip--

    // Create a connection to the gRPC server.
    grpcConn := grpc.Dial(
        "127.0.0.1:9090", // Or your gRPC server address.
        grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
    )
    defer grpcConn.Close()

    // Broadcast the tx via gRPC. We create a new client for the Protobuf Tx
    // service.
    txClient := tx.NewServiceClient(grpcConn)
    // We then call the BroadcastTx method on this client.
    grpcRes, err := txClient.BroadcastTx(
        ctx,
        &tx.BroadcastTxRequest{
            Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
            TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
        },
    )
    if err != nil {
        return err
    }

    fmt.Println(grpcRes.TxResponse.Code) // Should be `0` if the tx is successful

    return nil
}
```
