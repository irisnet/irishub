package rpc

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/query"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ethermint "github.com/evmos/ethermint/types"

	"github.com/irisnet/irishub/v4/app/keepers"
)

var _ authtypes.QueryServer = authQueryServer{}

type authQueryServer struct {
	key storetypes.StoreKey
	cdc codec.Codec
	authtypes.QueryServer
}

// overrideAuthServices overrides auth query service
func overrideAuthServices(cdc codec.Codec, cfg module.Configurator, appKeepers keepers.AppKeepers) {
	k := appKeepers.AccountKeeper
	key := appKeepers.GetKey(authtypes.StoreKey)
	ss := appKeepers.GetSubspace(authtypes.ModuleName)

	types.RegisterMsgServer(cfg.MsgServer(), authkeeper.NewMsgServerImpl(k))
	types.RegisterQueryServer(cfg.QueryServer(), authQueryServer{key: key, cdc: cdc, QueryServer: authkeeper.NewQueryServer(k)})

	m := authkeeper.NewMigrator(k, cfg.QueryServer(), ss)
	if err := cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 3 to 4: %v", types.ModuleName, err))
	}
}

// Since: cosmos-sdk 0.43
// Accounts returns all the existing accounts
func (a authQueryServer) Accounts(c context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(a.key)
	accountsStore := prefix.NewStore(store, types.AddressStoreKeyPrefix)

	var accounts []*codectypes.Any
	pageRes, err := query.Paginate(accountsStore, req.Pagination, func(key, value []byte) error {
		account := a.decodeAccount(value)
		any, err := codectypes.NewAnyWithValue(account)
		if err != nil {
			return err
		}

		accounts = append(accounts, any)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}

	return &types.QueryAccountsResponse{Accounts: accounts, Pagination: pageRes}, err
}

// Account returns account details based on address.
func (a authQueryServer) Account(c context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	res, err := a.AccountInfo(c, &authtypes.QueryAccountInfoRequest{
		Address: req.Address,
	})
	if err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(res.Info)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &types.QueryAccountResponse{Account: any}, nil
}

func (a authQueryServer) decodeAccount(bz []byte) sdk.AccountI {
	var acc sdk.AccountI
	if err := a.cdc.UnmarshalInterface(bz, &acc); err != nil {
		panic(err)
	}
	ethAcc, ok := acc.(*ethermint.EthAccount)
	if ok {
		return ethAcc.BaseAccount
	}
	return acc
}
