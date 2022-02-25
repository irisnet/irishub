# Messages

## MsgIssueDenom

This message defines a type of mt, there can be multiple mt of the same type

| **Field**        | **Type** | **Description**                                                                                                                        |
|:-----------------|:---------|:---------------------------------------------------------------------------------------------------------------------------------------|
| Name             | `string` | The name of the asset class.                       |
| Data             | `bytes`  | data is the app specific metadata of the MT class. Optional                                                                           |                                                                             |
| Sender           | `string` | The account address of the user sending the MT. By default it is __not__ required that the sender is also the owner of the MT.       |

```go
// MsgIssueDenom defines an SDK message for creating a new denom.
type MsgIssueDenom struct {
    Name   string
    Data   []byte
    Sender string
}
```

## MsgTransferDenom

This message is used by the owner of the MT classification to transfer the ownership of the MT classification to others

| **Field** | **Type** | **Description**                                                                         |
|:----------|:---------|:----------------------------------------------------------------------------------------|
| ID        | `string` | The unique ID of the Denom being transferred.                                           | 
| Sender    | `string` | The account address of the user sending the Denom.                                      |
| Recipient | `string` | The account address who will receive the Denom as a result of the transfer transaction. |

```go
// MsgTransferDenom defines an SDK message for transferring an denom to recipient.
type MsgTransferDenom struct {
    Id        string
    Sender    string
    Recipient string
}
```

## MsgMintMT

This message type is used for minting new tokens. If a new `MT` is minted under a new `Denom`, a new `Collection` will also be created, otherwise the `MT` is added to the existing `Collection`. If a new `MT` is minted by a new account, a new `Owner` is created, otherwise the `MT` `ID` is added to the existing `Owner`'s `IDCollection`. By default, anyone can execute this Message type. **It is highly recommended that a custom handler is made to restrict use of this Message type to prevent unintended use.**

| **Field** | **Type** | **Description**                                                                            |
|:----------|:---------|:-------------------------------------------------------------------------------------------|
| Id        | `string` | The unique ID of the MT being minted.                                                      |
| DenomId   | `string` | The unique ID of the denomination.                                                         |
| Amount    | `uint64` | The Amount of the MT being minted.                                                         |
| Data      | `bytes`  | The data of the MT.                                                                        |
| Sender    | `string` | The sender of the Message.                                                                 |
| Recipient | `string` | The recipient of the new MT.                                                               |

```go
// MsgMintMT defines an SDK message for creating a new MT.
type MsgMintMT struct {
    Id        string
	DenomId   string
    Amount    uint64 
    Data      []byte
    Sender    string
    Recipient string
}
```

## MsgEditMT

This message type allows the `Data` to be updated. **It is highly recommended that a custom handler is made to restrict use of this Message type to prevent unintended use.**

| **Field** | **Type** | **Description**                                                                                                  |
|:----------|:---------|:-----------------------------------------------------------------------------------------------------------------|
| Id        | `string` | The unique ID of the MT being edited.                                                                            |
| DenomId   | `string` | The unique ID of the denomination, necessary as multiple denominations are able to be represented on each chain. |
| Data      | `bytes`  | The data of the MT.                                                                                              |
| Sender    | `string` | The creator of the message.                                                                                      |

```go
// MsgEditMT defines an SDK message for editing an MT.
type MsgEditMT struct {
    Id      string
    DenomId string
    Data    []byte
    Sender  string
}
```

## MsgTransferMT

This is the most commonly expected MsgType to be supported across chains. While each application specific blockchain will have very different adoption of the `MsgMintMT`, `MsgBurnMT` and `MsgEditMT` it should be expected that most chains support the ability to transfer ownership of the MT asset. The exception to this would be non-transferable MTs that might be attached to reputation or some asset which should not be transferable. It still makes sense for this to be represented as an MT because there are common queries which will remain relevant to the MT type even if non-transferable. This Message will fail if the MT does not exist.

| **Field** | **Type** | **Description**                                                                                                                  |
|:----------|:---------|:---------------------------------------------------------------------------------------------------------------------------------|
| ID        | `string` | The unique ID of the MT being transferred.                                                                                       |
| DenomId   | `string` | The unique ID of the denomination, necessary as multiple denominations are able to be represented on each chain.                 |
| Amount    | `uint64` | The Amount of the MT being transferred.                                                                                          |
| Data      | `string` | The data of the MT.                                                                                                              |
| Sender    | `string` | The account address of the user sending the MT. By default it is __not__ required that the sender is also the owner of the MT.   |
| Recipient | `string` | The account address who will receive the MT as a result of the transfer transaction.                                             |

```go
// MsgTransferMT defines an SDK message for transferring an MT to recipient.
type MsgTransferMT struct {
    Id        string
    DenomId   string
    Amount    uint64
    Sender    string
    Recipient string
}
```

### MsgBurnMT

This message type is used for burning tokens which destroys.

| **Field** | **Type** | **Description**                                    |
|:----------|:---------|:---------------------------------------------------|
| Id        | `string` | The ID of the Token.                               |
| DenomId   | `string` | The Denom ID of the Token.                         |
| Amount    | `uint64` | The Amount of the MT being burning.                |
| Sender    | `string` | The account address of the user burning the token. |

```go
// MsgBurnMT defines an SDK message for burning an MT.
type MsgBurnMT struct {
    Id      string
    DenomId string
    Amount  uint64
    Sender  string
}
```

