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
	cdc codec.Codec
	key storetypes.StoreKey
	authkeeper.AccountKeeper
}

func (a authQueryServer) AccountAddressByID(ctx context.Context, request *authtypes.QueryAccountAddressByIDRequest) (*authtypes.QueryAccountAddressByIDResponse, error) {
	return &authtypes.QueryAccountAddressByIDResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) Params(ctx context.Context, request *authtypes.QueryParamsRequest) (*authtypes.QueryParamsResponse, error) {
	return &authtypes.QueryParamsResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) ModuleAccounts(ctx context.Context, request *authtypes.QueryModuleAccountsRequest) (*authtypes.QueryModuleAccountsResponse, error) {
	return &authtypes.QueryModuleAccountsResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) ModuleAccountByName(ctx context.Context, request *authtypes.QueryModuleAccountByNameRequest) (*authtypes.QueryModuleAccountByNameResponse, error) {
	return &authtypes.QueryModuleAccountByNameResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) Bech32Prefix(ctx context.Context, request *authtypes.Bech32PrefixRequest) (*authtypes.Bech32PrefixResponse, error) {
	return &authtypes.Bech32PrefixResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) AddressBytesToString(ctx context.Context, request *authtypes.AddressBytesToStringRequest) (*authtypes.AddressBytesToStringResponse, error) {
	return &authtypes.AddressBytesToStringResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) AddressStringToBytes(ctx context.Context, request *authtypes.AddressStringToBytesRequest) (*authtypes.AddressStringToBytesResponse, error) {
	return &authtypes.AddressStringToBytesResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

func (a authQueryServer) AccountInfo(ctx context.Context, request *authtypes.QueryAccountInfoRequest) (*authtypes.QueryAccountInfoResponse, error) {
	return &authtypes.QueryAccountInfoResponse{}, status.Error(codes.Unimplemented, "not implemented")
}

// overrideAuthServices overrides auth query service
func overrideAuthServices(cdc codec.Codec, cfg module.Configurator, appKeepers keepers.AppKeepers) {
	k := appKeepers.AccountKeeper
	key := appKeepers.GetKey(authtypes.StoreKey)
	ss := appKeepers.GetSubspace(authtypes.ModuleName)

	types.RegisterMsgServer(cfg.MsgServer(), authkeeper.NewMsgServerImpl(k))
	types.RegisterQueryServer(cfg.QueryServer(), authQueryServer{cdc, key, k})

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
	var acc types.AccountI
	if err := a.cdc.UnmarshalInterface(bz, &acc); err != nil {
		panic(err)
	}

	ethAcc, ok := acc.(*ethermint.EthAccount)
	if ok {
		return ethAcc.BaseAccount
	}
	return acc
}
