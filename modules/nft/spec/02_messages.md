# Messages

## MsgIssueDenom

This message defines a type of nft, there can be multiple nft of the
same type

| **Field** | **Type** | **Description**                                                                                                                  |
|:----------|:---------|:---------------------------------------------------------------------------------------------------------------------------------|
| Name      | `Id`     | The denomination ID of the NFT, necessary as multiple denominations are able to be represented on each chain.                    |
| Name      | `string` | The denomination name of the NFT, necessary as multiple denominations are able to be represented on each chain.                  |
| Sender    | `string` | The account address of the user sending the NFT. By default it is __not__ required that the sender is also the owner of the NFT. |
| Schema    | `string` | NFT specifications defined under this category                                                                                   |

```go
type MsgIssueDenom struct {
	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Schema string `protobuf:"bytes,3,opt,name=schema,proto3" json:"schema,omitempty"`
	Sender string `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
}
```

## MsgTransferNFT

This is the most commonly expected MsgType to be supported across
chains. While each application specific blockchain will have very
different adoption of the `MsgMintNFT`, `MsgBurnNFT` and `MsgEditNFT` it
should be expected that most chains support the ability to transfer
ownership of the NFT asset. The exception to this would be
non-transferable NFTs that might be attached to reputation or some asset
which should not be transferable. It still makes sense for this to be
represented as an NFT because there are common queriers which will
remain relevant to the NFT type even if non-transferable. This Message
will fail if the NFT does not exist. By default it will not fail if the
transfer is executed by someone beside the owner. **It is highly
recommended that a custom handler is made to restrict use of this
Message type to prevent unintended use.**

| **Field** | **Type** | **Description**                                                                                                                  |
|:----------|:---------|:---------------------------------------------------------------------------------------------------------------------------------|
| ID        | `string` | The unique ID of the NFT being transferred.                                                                                      |
| DenomId   | `string` | The unique ID of the denomination, necessary as multiple denominations are able to be represented on each chain.                 |
| Name      | `string` | The name of the NFT being transferred.                                                                                           |
| URI       | `string` | The URI pointing to a JSON object that contains subsequent tokenData information off-chain                                       |
| Data      | `string` | The data of the NFT                                                                                                              |
| Sender    | `string` | The account address of the user sending the NFT. By default it is __not__ required that the sender is also the owner of the NFT. |
| Recipient | `string` | The account address who will receive the NFT as a result of the transfer transaction.                                            |

```go
// MsgTransferNFT defines an SDK message for transferring an NFT to recipient.
type MsgTransferNFT struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DenomId   string `protobuf:"bytes,2,opt,name=denom_id,json=denomId,proto3" json:"denom_id,omitempty" yaml:"denom_id"`
	Name      string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	URI       string `protobuf:"bytes,4,opt,name=uri,proto3" json:"uri,omitempty"`
	Data      string `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Sender    string `protobuf:"bytes,6,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipient string `protobuf:"bytes,7,opt,name=recipient,proto3" json:"recipient,omitempty"`
}
```

## MsgEditNFT

This message type allows the `TokenURI` to be updated. By default anyone
can execute this Message type. **It is highly recommended that a custom
handler is made to restrict use of this Message type to prevent
unintended use.**

| **Field** | **Type** | **Description**                                                                                                  |
|:----------|:---------|:-----------------------------------------------------------------------------------------------------------------|
| Id        | `string` | The unique ID of the NFT being edited.                                                                           |
| DenomId   | `string` | The unique ID of the denomination, necessary as multiple denominations are able to be represented on each chain. |
| Name      | `string` | The name of the NFT being edited.                                                                                |
| URI       | `string` | The URI pointing to a JSON object that contains subsequent tokenData information off-chain                       |
| Data      | `string` | The data of the NFT                                                                                              |
| Sender    | `string` | The creator of the message                                                                                       |

```go
// MsgEditNFT defines an SDK message for editing a nft.
type MsgEditNFT struct {
	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DenomId string `protobuf:"bytes,2,opt,name=denom_id,json=denomId,proto3" json:"denom_id,omitempty" yaml:"denom_id"`
	Name    string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	URI     string `protobuf:"bytes,4,opt,name=uri,proto3" json:"uri,omitempty"`
	Data    string `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Sender  string `protobuf:"bytes,6,opt,name=sender,proto3" json:"sender,omitempty"`
}
```

## MsgMintNFT

This message type is used for minting new tokens. If a new `NFT` is
minted under a new `Denom`, a new `Collection` will also be created,
otherwise the `NFT` is added to the existing `Collection`. If a new
`NFT` is minted by a new account, a new `Owner` is created, otherwise
the `NFT` `ID` is added to the existing `Owner`'s `IDCollection`. By
default anyone can execute this Message type. **It is highly recommended
that a custom handler is made to restrict use of this Message type to
prevent unintended use.**

| **Field** | **Type** | **Description**                                                                            |
|:----------|:---------|:-------------------------------------------------------------------------------------------|
| ID        | `string` | The unique ID of the NFT being minted                                                      |
| DenomId   | `string` | The unique ID of the denomination.                                                         |
| Name      | `string` | The name of the NFT being minted.                                                          |
| URI       | `string` | The URI pointing to a JSON object that contains subsequent tokenData information off-chain |
| Data      | `string` | The data of the NFT.                                                                       |
| Sender    | `string` | The sender of the Message                                                                  |
| Recipient | `string` | The recipiet of the new NFT                                                                |

```go
// MsgMintNFT defines an SDK message for creating a new NFT.
type MsgMintNFT struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DenomId   string `protobuf:"bytes,2,opt,name=denom_id,json=denomId,proto3" json:"denom_id,omitempty" yaml:"denom_id"`
	Name      string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	URI       string `protobuf:"bytes,4,opt,name=uri,proto3" json:"uri,omitempty"`
	Data      string `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Sender    string `protobuf:"bytes,6,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipient string `protobuf:"bytes,7,opt,name=recipient,proto3" json:"recipient,omitempty"`
}
```

### MsgBurnNFT

This message type is used for burning tokens which destroys and deletes
them. By default anyone can execute this Message type. **It is highly
recommended that a custom handler is made to restrict use of this
Message type to prevent unintended use.**


| **Field** | **Type** | **Description**                                    |
|:----------|:---------|:---------------------------------------------------|
| ID        | `string` | The ID of the Token.                               |
| DenomId   | `string` | The Denom of the Token.                            |
| Sender    | `string` | The account address of the user burning the token. |

```go
// MsgBurnNFT defines an SDK message for burning a NFT.
type MsgBurnNFT struct {
	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DenomId string `protobuf:"bytes,2,opt,name=denom_id,json=denomId,proto3" json:"denom_id,omitempty" yaml:"denom_id"`
	Sender  string `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
}
```

