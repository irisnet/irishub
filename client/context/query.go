package context

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmliteErr "github.com/tendermint/tendermint/lite/errors"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GetNode returns an RPC client. If the context's client is not defined, an
// error is returned.
func (cliCtx CLIContext) GetNode() (rpcclient.Client, error) {
	if cliCtx.Client == nil {
		return nil, errors.New("no RPC client defined")
	}

	return cliCtx.Client, nil
}

// Query performs a query for information about the connected node.
func (cliCtx CLIContext) Query(path string, data cmn.HexBytes) (res []byte, err error) {
	return cliCtx.query(path, data)
}

// Query information about the connected node with a data payload
func (cliCtx CLIContext) QueryWithData(path string, data []byte) (res []byte, err error) {
	return cliCtx.query(path, data)
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (cliCtx CLIContext) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return cliCtx.queryStore(key, storeName, "key")
}

// QuerySubspace performs a query from a Tendermint node with the provided
// store name and subspace.
func (cliCtx CLIContext) QuerySubspace(subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, err := cliCtx.queryStore(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}

	cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(resRaw, &res)
	return
}

// GetAccount queries for an account given an address and a block height. An
// error is returned if the query or decoding fails.
func (cliCtx CLIContext) GetAccount(address []byte) (account auth.BaseAccount, err error) {
	if cliCtx.AccDecoder == nil {
		return account, errors.New("account decoder required but not provided")
	}

	res, err := cliCtx.queryAccount(address)
	if err != nil {
		return account, err
	}

	if err := cliCtx.Codec.UnmarshalJSON(res, &account); err != nil {
		return account, err
	}

	return account, nil
}

// GetFromAddress returns the from address from the context's name.
func (cliCtx CLIContext) GetFromAddress() (sdk.AccAddress, error) {
	return cliCtx.fromAddress, nil
}

// GetFromName returns the key name for the current context.
func (cliCtx CLIContext) GetFromName() (string, error) {
	return cliCtx.fromName, nil
}

// GetAccountNumber returns the next account number for the given account
// address.
func (cliCtx CLIContext) GetAccountNumber(address []byte) (uint64, error) {
	account, err := cliCtx.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return account.GetAccountNumber(), nil
}

// GetAccountSequence returns the sequence number for the given account
// address.
func (cliCtx CLIContext) GetAccountSequence(address []byte) (uint64, error) {
	account, err := cliCtx.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return account.GetSequence(), nil
}

// EnsureAccountExists ensures that an account exists for a given context. An
// error is returned if it does not.
func (cliCtx CLIContext) EnsureAccountExists() error {
	addr, err := cliCtx.GetFromAddress()
	if err != nil {
		return err
	}

	accountBytes, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), cliCtx.AccountStore)
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
func (cliCtx CLIContext) EnsureAccountExistsFromAddr(addr sdk.AccAddress) error {
	_, err := cliCtx.queryAccount(addr)
	return err
}

// queryAccount queries an account using custom query endpoint of auth module
// returns an error if result is `null` otherwise account data
func (cliCtx CLIContext) queryAccount(addr sdk.AccAddress) ([]byte, error) {
	bz, err := cliCtx.Codec.MarshalJSON(bank.NewQueryAccountParams(addr))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, bank.QueryAccount)

	res, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// query performs a query from a Tendermint node with the provided store name
// and path.
func (cliCtx CLIContext) query(path string, key cmn.HexBytes) (res []byte, err error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height: cliCtx.Height,
		Prove:  !cliCtx.TrustNode,
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
	if cliCtx.TrustNode || !isQueryStoreWithProof(path) {
		return resp.Value, nil
	}

	err = cliCtx.verifyProof(path, resp)
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// Verify verifies the consensus proof at given height.
func (cliCtx CLIContext) Verify(height int64) (tmtypes.SignedHeader, error) {
	check, err := tmliteProxy.GetCertifiedCommit(height, cliCtx.Client, cliCtx.Verifier)
	switch {
	case tmliteErr.IsErrCommitNotFound(err):
		return tmtypes.SignedHeader{}, ErrVerifyCommit(height)
	case err != nil:
		return tmtypes.SignedHeader{}, err
	}

	return check, nil
}

// verifyProof perform response proof verification.
func (cliCtx CLIContext) verifyProof(queryPath string, resp abci.ResponseQuery) error {
	if cliCtx.Verifier == nil {
		return fmt.Errorf("missing valid certifier to verify data from distrusted node")
	}

	// the AppHash for height H is in header H+1
	commit, err := cliCtx.Verify(resp.Height + 1)
	if err != nil {
		return err
	}

	// TODO: Instead of reconstructing, stash on CLIContext field?
	prt := store.DefaultProofRuntime()

	// TODO: Better convention for path?
	storeName, err := parseQueryStorePath(queryPath)
	if err != nil {
		return err
	}

	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte(storeName), merkle.KeyEncodingURL)
	kp = kp.AppendKey(resp.Key, merkle.KeyEncodingURL)

	if resp.Value == nil {
		err = prt.VerifyAbsence(resp.Proof, commit.Header.AppHash, kp.String())
		if err != nil {
			return errors.Wrap(err, "failed to prove merkle proof")
		}
		return nil
	}
	err = prt.VerifyValue(resp.Proof, commit.Header.AppHash, kp.String(), resp.Value)
	if err != nil {
		return errors.Wrap(err, "failed to prove merkle proof")
	}

	return nil
}

// queryStore performs a query from a Tendermint node with the provided a store
// name and path.
func (cliCtx CLIContext) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return cliCtx.query(path, key)
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType can be app or store.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}
	paths := strings.SplitN(path[1:], "/", 3)
	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case store.RequireProof("/" + paths[2]):
		return true
	}
	return false
}

// parseQueryStorePath expects a format like /store/<storeName>/key.
func parseQueryStorePath(path string) (storeName string, err error) {
	if !strings.HasPrefix(path, "/") {
		return "", errors.New("expected path to start with /")
	}
	paths := strings.SplitN(path[1:], "/", 3)
	switch {
	case len(paths) != 3:
		return "", errors.New("expected format like /store/<storeName>/key")
	case paths[0] != "store":
		return "", errors.New("expected format like /store/<storeName>/key")
	case paths[2] != "key":
		return "", errors.New("expected format like /store/<storeName>/key")
	}
	return paths[1], nil
}

func (cliCtx CLIContext) GetCoinType(coinName string) (sdk.CoinType, error) {
	var coinType sdk.CoinType
	coinName = strings.ToLower(coinName)
	if coinName == "" {
		return sdk.CoinType{}, fmt.Errorf("coin name is empty")
	}
	if coinName == sdk.Iris {
		coinType = sdk.IrisCoinType
	} else {
		params := asset.QueryTokenParams{
			TokenId: coinName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			return sdk.CoinType{}, err
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryToken), bz)
		if err != nil {
			return sdk.CoinType{}, fmt.Errorf("unsupported coin type \"%s\"", coinName)
		}

		var token asset.FungibleToken
		err = cliCtx.Codec.UnmarshalJSON(res, &token)
		if err != nil {
			return sdk.CoinType{}, err
		}

		coinType = token.GetCoinType()
	}

	return coinType, nil
}

func (cliCtx CLIContext) ConvertToMainUnit(coinsStr string) (coins []string, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return coins, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		mainUnit, err := sdk.GetCoinName(coinStr)
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
	mainUnit, err := sdk.GetCoinName(coinStr)
	coinType, err := cliCtx.GetCoinType(mainUnit)
	if err != nil {
		return sdk.Coin{}, err
	}

	coin, err := coinType.ConvertToMinDenomCoin(coinStr)
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
	coinMap := make(map[string]sdk.Coin)
	for _, coinStr := range coinStrs {
		coin, err := cliCtx.ParseCoin(coinStr)
		if err != nil {
			return sdk.Coins{}, err
		}
		if _, ok := coinMap[coin.Denom]; ok {
			coinMap[coin.Denom] = coinMap[coin.Denom].Plus(coin)
		} else {
			coinMap[coin.Denom] = coin
		}
	}

	for _, coin := range coinMap {
		coins = append(coins, coin)
	}
	coins = coins.Sort()
	return coins, nil
}

func (cliCtx CLIContext) NetInfo() (*ctypes.ResultNetInfo, error) {
	client, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}
	httpClient := client.(*tmclient.HTTP)
	return httpClient.NetInfo()
}

func (cliCtx CLIContext) GetLatestHeight() (int64, error) {
	client, err := cliCtx.GetNode()
	if err != nil {
		return 0, err
	}
	httpClient := client.(*tmclient.HTTP)

	status, err := httpClient.Status()
	if err != nil {
		return 0, err
	}
	if status.SyncInfo.CatchingUp {
		return 0, fmt.Errorf("the connected full node is still syncing blocks")
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (cliCtx CLIContext) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	client := &http.Client{}
	url := strings.Replace(cliCtx.NodeURI, "tcp", "http", 1)
	reqUri := fmt.Sprintf("%s/%s", url, "num_unconfirmed_txs")

	resp, err := client.Get(reqUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res = struct {
		JsonRpc string                      `json:"jsonrpc"`
		Id      string                      `json:"id"`
		Result  ctypes.ResultUnconfirmedTxs `json:"result"`
	}{}

	if err := cliCtx.Codec.UnmarshalJSON(body, &res); err != nil {
		return nil, err
	}

	return &res.Result, nil
}

// PrintOutput prints output while respecting output and indent flags
// NOTE: pass in marshalled structs that have been unmarshaled
// because this function will panic on marshaling errors
func (ctx CLIContext) PrintOutput(toPrint fmt.Stringer) (err error) {
	var out []byte

	switch ctx.OutputFormat {
	case "text":
		humanStringer, ok := toPrint.(sdk.Stringer)
		if ok {
			out = []byte(humanStringer.HumanString(ctx))
		} else {
			out = []byte(toPrint.String())
		}

	case "json":
		if ctx.Indent {
			out, err = ctx.Codec.MarshalJSONIndent(toPrint, "", "  ")
		} else {
			out, err = ctx.Codec.MarshalJSON(toPrint)
		}
	}

	if err != nil {
		return
	}

	fmt.Println(string(out))
	return
}
