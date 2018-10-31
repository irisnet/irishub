package context

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/pkg/errors"

	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmliteErr "github.com/tendermint/tendermint/lite/errors"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/app"
)

// GetNode returns an RPC client. If the context's client is not defined, an
// error is returned.
func (ctx CLIContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("no RPC client defined")
	}

	return ctx.Client, nil
}

// Query performs a query for information about the connected node.
func (ctx CLIContext) Query(path string, data cmn.HexBytes) (res []byte, err error) {
	return ctx.query(path, data)
}

// Query information about the connected node with a data payload
func (ctx CLIContext) QueryWithData(path string, data []byte) (res []byte, err error) {
	return ctx.query(path, data)
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx CLIContext) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.queryStore(key, storeName, "key")
}

// QuerySubspace performs a query from a Tendermint node with the provided
// store name and subspace.
func (ctx CLIContext) QuerySubspace(subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, err := ctx.queryStore(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}

	ctx.Codec.MustUnmarshalBinary(resRaw, &res)
	return
}

// GetAccount queries for an account given an address and a block height. An
// error is returned if the query or decoding fails.
func (ctx CLIContext) GetAccount(address []byte) (auth.Account, error) {
	if ctx.AccDecoder == nil {
		return nil, errors.New("account decoder required but not provided")
	}

	res, err := ctx.QueryStore(auth.AddressStoreKey(address), ctx.AccountStore)
	if err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, err
	}

	account, err := ctx.AccDecoder(res)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetFromAddress returns the from address from the context's name.
func (ctx CLIContext) GetFromAddress() (sdk.AccAddress, error) {
	return ctx.fromAddress, nil
}

// GetFromName returns the key name for the current context.
func (ctx CLIContext) GetFromName() (string, error) {
	return ctx.fromName, nil
}

// GetAccountNumber returns the next account number for the given account
// address.
func (ctx CLIContext) GetAccountNumber(address []byte) (int64, error) {
	account, err := ctx.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return account.GetAccountNumber(), nil
}

// GetAccountSequence returns the sequence number for the given account
// address.
func (ctx CLIContext) GetAccountSequence(address []byte) (int64, error) {
	account, err := ctx.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return account.GetSequence(), nil
}

// EnsureAccountExists ensures that an account exists for a given context. An
// error is returned if it does not.
func (ctx CLIContext) EnsureAccountExists() error {
	addr, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	accountBytes, err := ctx.QueryStore(auth.AddressStoreKey(addr), ctx.AccountStore)
	if err != nil {
		return err
	}

	if len(accountBytes) == 0 {
		return ErrInvalidAccount(addr)
	}

	return nil
}

// EnsureAccountExistsFromAddr ensures that an account exists for a given
// address. Instead of using the context's from name, a direct address is
// given. An error is returned if it does not.
func (ctx CLIContext) EnsureAccountExistsFromAddr(addr sdk.AccAddress) error {
	accountBytes, err := ctx.QueryStore(auth.AddressStoreKey(addr), ctx.AccountStore)
	if err != nil {
		return err
	}

	if len(accountBytes) == 0 {
		return ErrInvalidAccount(addr)
	}

	return nil
}

// query performs a query from a Tendermint node with the provided store name
// and path.
func (ctx CLIContext) query(path string, key cmn.HexBytes) (res []byte, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height:  ctx.Height,
		Trusted: ctx.TrustNode,
	}

	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.Errorf(resp.Log)
	}

	// data from trusted node or subspace query doesn't need verification
	if ctx.TrustNode || !isQueryStoreWithProof(path) {
		return resp.Value, nil
	}

	err = ctx.verifyProof(path, resp)
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// Verify verifies the consensus proof at given height.
func (ctx CLIContext) Verify(height int64) (tmtypes.SignedHeader, error) {
	check, err := tmliteProxy.GetCertifiedCommit(height, ctx.Client, ctx.Verifier)
	switch {
	case tmliteErr.IsErrCommitNotFound(err):
		return tmtypes.SignedHeader{}, ErrVerifyCommit(height)
	case err != nil:
		return tmtypes.SignedHeader{}, err
	}

	return check, nil
}

// verifyProof perform response proof verification.
func (ctx CLIContext) verifyProof(_ string, resp abci.ResponseQuery) error {
	if ctx.Verifier == nil {
		return fmt.Errorf("missing valid certifier to verify data from distrusted node")
	}

	// the AppHash for height H is in header H+1
	commit, err := ctx.Verify(resp.Height + 1)
	if err != nil {
		return err
	}

	var multiStoreProof store.MultiStoreProof
	cdc := codec.New()

	err = cdc.UnmarshalBinary(resp.Proof, &multiStoreProof)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshalBinary rangeProof")
	}

	// verify the substore commit hash against trusted appHash
	substoreCommitHash, err := store.VerifyMultiStoreCommitInfo(
		multiStoreProof.StoreName, multiStoreProof.StoreInfos, commit.Header.AppHash,
	)
	if err != nil {
		return errors.Wrap(err, "failed in verifying the proof against appHash")
	}

	err = store.VerifyRangeProof(resp.Key, resp.Value, substoreCommitHash, multiStoreProof.RangeProof)
	if err != nil {
		return errors.Wrap(err, "failed in the range proof verification")
	}

	return nil
}

// queryStore performs a query from a Tendermint node with the provided a store
// name and path.
func (ctx CLIContext) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return ctx.query(path, key)
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType can be app or store.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}
	paths := strings.SplitN(path[1:], "/", 3)
	if len(paths) != 3 {
		return false
	}

	if store.RequireProof("/" + paths[2]) {
		return true
	}
	return false
}

func (cliCtx CLIContext) GetCoinType(coinName string) (types.CoinType, error) {
	var coinType types.CoinType
	coinName = strings.ToLower(coinName)
	if coinName == "" {
		return types.CoinType{}, fmt.Errorf("coin name is empty")
	}
	if coinName == app.Denom {
		coinType = app.IrisCt
	} else {
		key := types.CoinTypeKey(coinName)
		bz, err := cliCtx.QueryStore([]byte(key), "params")
		if err != nil {
			return coinType, err
		}

		if bz == nil {
			return types.CoinType{}, fmt.Errorf("unsupported coin type \"%s\"", coinName)
		}

		if err = cliCtx.Codec.UnmarshalBinary(bz, &coinType); err != nil {
			return coinType, err
		}
	}

	return coinType, nil
}

func (cliCtx CLIContext) ConvertCoinToMainUnit(coinsStr string) (coins []string, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return coins, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		mainUnit, err := types.GetCoinName(coinStr)
		coinType, err := cliCtx.GetCoinType(mainUnit)
		if err != nil {
			return nil, err
		}

		coin, err := coinType.Convert(coinStr, mainUnit)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

func (cliCtx CLIContext) ParseCoin(coinStr string) (sdk.Coin, error) {
	mainUnit, err := types.GetCoinName(coinStr)
	coinType, err := cliCtx.GetCoinType(mainUnit)
	if err != nil {
		return sdk.Coin{}, err
	}

	coin, err := coinType.ConvertToMinCoin(coinStr)
	if err != nil {
		return sdk.Coin{}, err
	}
	return coin, nil
}

func (cliCtx CLIContext) ParseCoins(coinsStr string) (coins sdk.Coins, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return coins, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := cliCtx.ParseCoin(coinStr)
		if err != nil {
			return coins, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}
