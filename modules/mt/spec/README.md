<!--
order: 0
title: MT Overview
parent:
  title: "MT"
-->

# MT Specification

## Overview

The MT Module described here is meant to be used as a module across chains for managing non-fungible token that represent individual assets with unique features. This standard was first developed on Ethereum within the ERC-721 and the subsequent EIP of the same name. This standard utilized the features of the Ethereum blockchain as well as the restrictions. The subsequent ERC-1155 standard addressed some of the restrictions of Ethereum regarding storage costs and semi-fungible assets.

MTs on application specific blockchains share some but not all features as their Ethereum brethren. Since application specific blockchains are more flexible in how their resources are utilized it makes sense that should have the option of exploiting those resources. This includes the aility to use strings as IDs and to optionally store tokenData on chain. The user-flow of composability with smart contracts should also be rethought on application specific blockchains with regard to Inter-Blockchain Communication as it is a different design experience from communication between smart contracts.

## Contents

1. **[State](./01_state.md)**
   - [MT](./01_state.md#mt)
   - [Collections](./01_state.md#collections)
   - [Owners](./01_state.md#owners)
1. **[Messages](./02_messages.md)**
   - [Issue Denom](./02_messages.md#msgissuedenom)
   - [Transfer Denom](./02_messages.md#msgtransferdenom)
   - [Transfer MT](./02_messages.md#msgtransfermt)
   - [Edit MT](./02_messages.md#msgtransfermt)
   - [Mint MT](./02_messages.md#msgmintmt)
   - [Burn MT](./02_messages.md#msgburnmt)
1. **[Events](./03_events.md)**
   - [Handlers](03_events.md#handlers)
1. **[Future Improvements](./04_future_improvements.md)**

## A Note on Metadata & IBC

The MT includes `tokenURI` in order to be backwards compatible with Ethereum based MTs. However the `MT` type is an interface that allows arbitrary tokenData to be stored on chain should it need be. Originally the module included `name`, `description` and `image` to demonstrate these capabilities. They were removed in order for the MT to be more efficient for use cases that don't include a need for that information to be stored on chain. A demonstration of including them will be included in a sample app. It is also under discussion to move all tokenData to a separate module that can handle arbitrary amounts of data on chain and can be used to describe assets beyond Non-Fungible Tokens, like normal Fungible Tokens `Coin` that could describe attributes like decimal places and vesting status.

A stand-alone tokenData Module would allow for independent standards to evolve regarding arbitrary asset types with expanding precision. The standards supported by [http://schema.org](http://schema.org) and the process of adding nested information is being considered as a starting point for that standard. The Blockchain Gaming Alliance is working on a tokenData standard to be used for specifically blockchain gaming assets.

With regards to Inter-Blockchain Communication the responsibility of the integrity of the tokenData should be left to the origin chain. If a secondary chain was responsible for storing the source of truth of the tokenData for an asset tracking that source of truth would become difficult if not impossible to track. Since origin chains are where the design and use of the MT is determined, it should be up to that origin chain to decide who can update tokenData and under what circumstances. Secondary chains can use IBC queriers to check needed tokenData or keep redundant copies of the tokenData locally when they receive the MT originally. In that case it should be up to te secondary chain to keep the tokenData in sync if need be, similar to how layer 2 solutions keep tokenData in sync with a source of truth using events.

## Custom App-Specific Handlers

Each message type comes with a default handler that can be used by default but will most likely be too limited for each use case. In order to make them useful for as many situations as possible, there are very few limitations on who can execute the Messages and do things like mint, burn or edit tokenData. We recommend that custom handlers are created to add in custom logic and restrictions over when the Message types can be executed. Below is an example implementation for initializing the module within the module manager so that a custom handler can be added. This can be seen in the example [MT app](https://github.com/okwme/cosmos-mt).

```go
// custom-handler.go

// OverrideMTModule overrides the MT module for custom handlers
type OverrideMTModule struct {
    mt.AppModule
    k mt.Keeper
}

// NewHandler overwrites the legacy NewHandler in order to allow custom logic for handling the messages
func (am OverrideMTModule) NewHandler() sdk.Handler {
    return CustomMTHandler(am.k)
}

// NewOverrideMTModule generates a new MT Module
func NewOverrideMTModule(appModule mt.AppModule, keeper mt.Keeper) OverrideMTModule {
    return OverrideMTModule{
        AppModule: appModule,
        k:         keeper,
    }
}
```

You can see here that `OverrideMTModule` is the same as `mt.AppModule` except for the `NewHandler()` method. This method now returns a new Handler called `CustomMTHandler`. This custom handler can be seen below:

```go
// CustomMTHandler routes the messages to the handlers
func CustomMTHandler(k keeper.Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
        switch msg := msg.(type) {
        case types.MsgTransferMT:
            return mt.HandleMsgTransferMT(ctx, msg, k)
        case types.MsgEditMT:
            return mt.HandleMsgEditdata(ctx, msg, k)
        case types.MsgMintMT:
            return HandleMsgMintMTCustom(ctx, msg, k) // <-- This one is custom, the others fall back onto the default
        case types.MsgBurnMT:
            return mt.HandleMsgBurnMT(ctx, msg, k)
        default:
            errMsg := fmt.Sprintf("unrecognized mt message type: %T", msg)
            return sdk.ErrUnknownRequest(errMsg).Result()
        }
    }
}

// HandleMsgMintMTCustom is a custom handler that handles MsgMintMT
func HandleMsgMintMTCustom(ctx sdk.Context, msg types.MsgMintMT, k keeper.Keeper) sdk.Result {
    isTwilight := checkTwilight(ctx)

    if isTwilight {
        return mt.HandleMsgMintMT(ctx, msg, k)
    }

    errMsg := fmt.Sprintf("Can't mint astral bodies outside of twilight!")
    return sdk.ErrUnknownRequest(errMsg).Result()
}
```

The default handlers are imported here with the MT module and used for `MsgTransferMT`, `MsgEditMT` and `MsgBurnMT`. The `MsgMintMT` however is handled with a custom function called `HandleMsgMintMTCustom`. This custom function also utilizes the imported MT module handler `HandleMsgMintMT`, but only after certain conditions are checked. In this case it checks a function called `checkTwilight` which returns a boolean. Only if `isTwilight` is true will the Message succeed.

This pattern of inheriting and utilizing the module handlers wrapped in custom logic should allow each application specific blockchain to use the MT while customizing it to their specific requirements.
