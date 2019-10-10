package operations

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/irisnet/irishub/modules/asset"
)

// SimulateMsgCreateGateway generates a MsgCreateGateway with random values.
func SimulateMsgCreateGateway(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		owner := simulation.RandomAcc(r, accs)
		moniker := simulation.RandStringOfLength(r, 5)
		identity := simulation.RandStringOfLength(r, 10)
		details := simulation.RandStringOfLength(r, 50)
		website := simulation.RandStringOfLength(r, 20)

		msg := asset.NewMsgCreateGateway(owner.Address, moniker, identity, details, website)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgEditGateway generates a MsgEditGateway with random values.
func SimulateMsgEditGateway(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		owner := simulation.RandomAcc(r, accs)
		moniker := simulation.RandStringOfLength(r, 5)
		identity := simulation.RandStringOfLength(r, 10)
		details := simulation.RandStringOfLength(r, 50)
		website := simulation.RandStringOfLength(r, 20)

		msg := asset.NewMsgEditGateway(owner.Address, moniker, identity, details, website)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgTransferGatewayOwner generates a MsgTransferGatewayOwner with random values.
func SimulateMsgTransferGatewayOwner(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		owner := simulation.RandomAcc(r, accs)
		moniker := simulation.RandStringOfLength(r, 5)
		to := simulation.RandomAcc(r, accs)

		msg := asset.NewMsgTransferGatewayOwner(owner.Address, moniker, to.Address)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgIssueToken generates a MsgIssueToken with random values.
func SimulateMsgIssueToken(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		owner := simulation.RandomAcc(r, accs)
		family := asset.AssetFamily(simulation.RandIntBetween(r, 0, 1))
		source := asset.AssetSource(simulation.RandIntBetween(r, 0, 2))
		gateway := simulation.RandStringOfLength(r, 5)
		symbol := simulation.RandStringOfLength(r, 5)
		canonicalSymbol := simulation.RandStringOfLength(r, 5)
		name := simulation.RandStringOfLength(r, 8)
		decimal := uint8(simulation.RandIntBetween(r, 8, 18))
		minUnitAlias := simulation.RandStringOfLength(r, 3)
		initialSupply := r.Uint64()
		maxSupply := initialSupply + uint64(simulation.RandIntBetween(r, 0, 1000))
		mintable, _ := asset.ParseBool(fmt.Sprintf("%d", simulation.RandIntBetween(r, 0, 1)))

		msg := asset.NewMsgIssueToken(family, source, gateway, symbol, canonicalSymbol, name, decimal, minUnitAlias, initialSupply, maxSupply, mintable, owner.Address)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgMintToken generates a MsgMintToken with random values.
func SimulateMsgMintToken(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		owner := simulation.RandomAcc(r, accs)
		to := simulation.RandomAcc(r, accs)
		tokenId := simulation.RandStringOfLength(r, 5)
		amount := r.Uint64()

		msg := asset.NewMsgMintToken(tokenId, owner.Address, to.Address, amount)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgEditToken generates a MsgEditToken with random values.
func SimulateMsgEditToken(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		owner := simulation.RandomAcc(r, accs)
		name := simulation.RandStringOfLength(r, 8)
		canonicalSymbol := simulation.RandStringOfLength(r, 5)
		minUnitAlias := simulation.RandStringOfLength(r, 3)
		tokenId := simulation.RandStringOfLength(r, 5)
		maxSupply := r.Uint64()
		mintable, _ := asset.ParseBool(fmt.Sprintf("%d", simulation.RandIntBetween(r, 0, 1)))

		msg := asset.NewMsgEditToken(name, canonicalSymbol, minUnitAlias, tokenId, maxSupply, mintable, owner.Address)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgTransferTokenOwner generates a MsgTransferTokenOwner with random values.
func SimulateMsgTransferTokenOwner(k asset.Keeper) simulation.Operation {
	handler := asset.NewHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		srcOwner := simulation.RandomAcc(r, accs)
		dstOwner := simulation.RandomAcc(r, accs)
		tokenId := simulation.RandStringOfLength(r, 5)

		msg := asset.NewMsgTransferTokenOwner(srcOwner.Address, dstOwner.Address, tokenId)

		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(asset.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := handler(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
