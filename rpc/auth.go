package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/query"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ethermint "github.com/evmos/ethermint/types"
)

var _ authtypes.QueryServer = authServer{}

type authServer struct {
	key storetypes.StoreKey
	k   authkeeper.AccountKeeper
}

func RegisterAuthServer(cfg module.Configurator, key storetypes.StoreKey, k authkeeper.AccountKeeper) {
	authtypes.RegisterQueryServer(cfg.QueryServer(), authServer{key, k})
	m := authkeeper.NewMigrator(k, cfg.QueryServer())
	err := cfg.RegisterMigration(authtypes.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(err)
	}

	err = cfg.RegisterMigration(authtypes.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(err)
	}
}

// Since: cosmos-sdk 0.43
// Accounts returns all the existing accounts
func (a authServer) Accounts(c context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
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
func (a authServer) Account(c context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "Address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	account := a.k.GetAccount(ctx, addr)
	if account == nil {
		return nil, status.Errorf(codes.NotFound, "account %s not found", req.Address)
	}

	ethAcc, ok := account.(*ethermint.EthAccount)
	if ok {
		account = ethAcc.BaseAccount
	}

	any, err := codectypes.NewAnyWithValue(account)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &types.QueryAccountResponse{Account: any}, nil
}

// AccountAddressByID returns account address based on account number.
//
// Since: cosmos-sdk 0.46.2
func (a authServer) AccountAddressByID(c context.Context, req *types.QueryAccountAddressByIDRequest) (*types.QueryAccountAddressByIDResponse, error) {
	return a.k.AccountAddressByID(c, req)
}

// Params queries all parameters.
func (a authServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return a.k.Params(ctx, req)
}

// ModuleAccounts returns all the existing module accounts.
//
// Since: cosmos-sdk 0.46
func (a authServer) ModuleAccounts(ctx context.Context, req *types.QueryModuleAccountsRequest) (*types.QueryModuleAccountsResponse, error) {
	return a.k.ModuleAccounts(ctx, req)
}

// ModuleAccountByName returns the module account info by module name
func (a authServer) ModuleAccountByName(ctx context.Context, req *types.QueryModuleAccountByNameRequest) (*types.QueryModuleAccountByNameResponse, error) {
	return a.k.ModuleAccountByName(ctx, req)
}

// Bech32Prefix queries bech32Prefix
//
// Since: cosmos-sdk 0.46
func (a authServer) Bech32Prefix(ctx context.Context, req *types.Bech32PrefixRequest) (*types.Bech32PrefixResponse, error) {
	return a.k.Bech32Prefix(ctx, req)
}

// AddressBytesToString converts Account Address bytes to string
//
// Since: cosmos-sdk 0.46
func (a authServer) AddressBytesToString(ctx context.Context, req *types.AddressBytesToStringRequest) (*types.AddressBytesToStringResponse, error) {
	return a.k.AddressBytesToString(ctx, req)
}

// AddressStringToBytes converts Address string to bytes
//
// Since: cosmos-sdk 0.46
func (a authServer) AddressStringToBytes(ctx context.Context, req *types.AddressStringToBytesRequest) (*types.AddressStringToBytesResponse, error) {
	return a.k.AddressStringToBytes(ctx, req)
}

func (a authServer) decodeAccount(bz []byte) types.AccountI {
	acc, err := a.k.UnmarshalAccount(bz)
	if err != nil {
		panic(err)
	}

	ethAcc, ok := acc.(*ethermint.EthAccount)
	if ok {
		return ethAcc.BaseAccount
	}
	return acc
}
