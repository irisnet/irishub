package context

import (
	"fmt"
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/lite"
	tmliteErr "github.com/tendermint/tendermint/lite/errors"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
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
func (cliCtx CLIContext) Query(path string) (res []byte, err error) {
	return cliCtx.query(path, nil)
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

	cliCtx.Codec.MustUnmarshalBinary(resRaw, &res)
	return
}

// GetAccount queries for an account given an address and a block height. An
// error is returned if the query or decoding fails.
func (cliCtx CLIContext) GetAccount(address []byte) (auth.Account, error) {
	if cliCtx.AccDecoder == nil {
		return nil, errors.New("account decoder required but not provided")
	}

	res, err := cliCtx.QueryStore(auth.AddressStoreKey(address), cliCtx.AccountStore)
	if err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, err
	}

	account, err := cliCtx.AccDecoder(res)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetFromAddress returns the from address from the context's name.
func (cliCtx CLIContext) GetFromAddress() (from sdk.AccAddress, err error) {
	kb, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}
	if cliCtx.GenerateOnly {
		signerAddress, err := sdk.AccAddressFromBech32(cliCtx.SignerAddr)
		// When generate-only is true, if the user specified signer address is correct, then just return the address.
		if err == nil {
			return signerAddress, nil
		}
	}
	if cliCtx.FromAddressName == "" {
		return nil, fmt.Errorf("must provide a from address name")
	}
	info, err := kb.Get(cliCtx.FromAddressName)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(info.GetPubKey().Address()), nil
}

// GetAccountNumber returns the next account number for the given account
// address.
func (cliCtx CLIContext) GetAccountNumber(address []byte) (int64, error) {
	account, err := cliCtx.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return account.GetAccountNumber(), nil
}

// GetAccountSequence returns the sequence number for the given account
// address.
func (cliCtx CLIContext) GetAccountSequence(address []byte) (int64, error) {
	account, err := cliCtx.GetAccount(address)
	if err != nil {
		return 0, err
	}

	return account.GetSequence(), nil
}

// BroadcastTx broadcasts transaction bytes to a Tendermint node.
func (cliCtx CLIContext) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if !res.CheckTx.IsOK() {
		return res, errors.Errorf("checkTx failed: (%d) %s",
			res.CheckTx.Code,
			res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return res, errors.Errorf("deliverTx failed: (%d) %s",
			res.DeliverTx.Code,
			res.DeliverTx.Log)
	}

	return res, err
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (cliCtx CLIContext) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxAsync(tx)
	if err != nil {
		return res, err
	}

	return res, err
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
	accountBytes, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), cliCtx.AccountStore)
	if err != nil {
		return err
	}

	if len(accountBytes) == 0 {
		return ErrInvalidAccount(addr)
	}

	return nil
}

// EnsureBroadcastTx broadcasts a transactions either synchronously or
// asynchronously based on the context parameters. The result of the broadcast
// is parsed into an intermediate structure which is logged if the context has
// a logger defined.
func (cliCtx CLIContext) EnsureBroadcastTx(txBytes []byte) error {
	if cliCtx.Async {
		return cliCtx.ensureBroadcastTxAsync(txBytes)
	}

	return cliCtx.ensureBroadcastTx(txBytes)
}

func (cliCtx CLIContext) ensureBroadcastTxAsync(txBytes []byte) error {
	res, err := cliCtx.BroadcastTxAsync(txBytes)
	if err != nil {
		return err
	}

	if cliCtx.JSON {
		type toJSON struct {
			TxHash string
		}

		if cliCtx.Logger != nil {
			resJSON := toJSON{res.Hash.String()}
			bz, err := cliCtx.Codec.MarshalJSON(resJSON)
			if err != nil {
				return err
			}

			cliCtx.Logger.Write(bz)
			io.WriteString(cliCtx.Logger, "\n")
		}
	} else {
		if cliCtx.Logger != nil {
			io.WriteString(cliCtx.Logger, fmt.Sprintf("Async tx sent (tx hash: %s)\n", res.Hash))
		}
	}

	return nil
}

func (cliCtx CLIContext) ensureBroadcastTx(txBytes []byte) error {
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	if cliCtx.JSON {
		// since JSON is intended for automated scripts, always include
		// response in JSON mode.
		type toJSON struct {
			Height   int64
			TxHash   string
			Response string
		}

		if cliCtx.Logger != nil {
			resJSON := toJSON{res.Height, res.Hash.String(), fmt.Sprintf("%+v", res.DeliverTx)}
			bz, err := cliCtx.Codec.MarshalJSON(resJSON)
			if err != nil {
				return err
			}

			cliCtx.Logger.Write(bz)
			io.WriteString(cliCtx.Logger, "\n")
		}

		return nil
	}

	if cliCtx.Logger != nil {
		resStr := fmt.Sprintf("Committed at block %d (tx hash: %s)\n", res.Height, res.Hash.String())

		if cliCtx.PrintResponse {
			jsonStr,_ := DeliverTxMarshalIndentJSON(res.DeliverTx)
			resStr = fmt.Sprintf("Committed at block %d (tx hash: %s, response: %+v)\n%s\n",
				res.Height, res.Hash.String(), res.DeliverTx,string(jsonStr),
			)

		}

		io.WriteString(cliCtx.Logger, resStr)
	}

	return nil
}

func DeliverTxMarshalIndentJSON(dtx abci.ResponseDeliverTx) ([]byte, error) {

	tags := make(map[string]string)
	for _, kv := range dtx.Tags {
		tags[string(kv.Key)] = strings.Replace(string(kv.Value),"\\","",-1)
	}

	return json.MarshalIndent(&struct {
		Tags map[string]string `json:"tags,omitempty"`
	}{
		Tags:      tags,
	}, " ", "  ")
}

// query performs a query from a Tendermint node with the provided store name
// and path.
func (cliCtx CLIContext) query(path string, key common.HexBytes) (res []byte, err error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height:  cliCtx.Height,
		Trusted: cliCtx.TrustNode,
	}

	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.Errorf("query failed: (%d) %s", resp.Code, resp.Log)
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

// verifyProof perform response proof verification
func (cliCtx CLIContext) verifyProof(path string, resp abci.ResponseQuery) error {

	if cliCtx.Certifier == nil {
		return fmt.Errorf("missing valid certifier to verify data from untrusted node")
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return err
	}

	// AppHash for height H is in header H+1
	commit, err := tmliteProxy.GetCertifiedCommit(resp.Height+1, node, cliCtx.Certifier)
	if err != nil {
		return err
	}

	var multiStoreProof store.MultiStoreProof
	cdc := wire.NewCodec()
	err = cdc.UnmarshalBinary(resp.Proof, &multiStoreProof)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshalBinary rangeProof")
	}

	// Verify the substore commit hash against trusted appHash
	substoreCommitHash, err := store.VerifyMultiStoreCommitInfo(
		multiStoreProof.StoreName, multiStoreProof.StoreInfos, commit.Header.AppHash)
	if err != nil {
		return errors.Wrap(err, "failed in verifying the proof against appHash")
	}

	err = store.VerifyRangeProof(resp.Key, resp.Value, substoreCommitHash, &multiStoreProof.RangeProof)
	if err != nil {
		return errors.Wrap(err, "failed in the range proof verification")
	}

	return nil
}

// queryStore performs a query from a Tendermint node with the provided a store
// name and path.
func (cliCtx CLIContext) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return cliCtx.query(path, key)
}

func (cliCtx CLIContext) GetCoinType(coinName string) (types.CoinType, error) {
	var coinType types.CoinType
	coinName = strings.ToLower(coinName)
	if coinName == app.Denom {
		coinType = app.IrisCt
	} else {
		key := types.CoinTypeKey(coinName)
		bz, err := cliCtx.QueryStore([]byte(key), "iparams")
		if err != nil {
			return coinType, err
		}

		if bz == nil {
			return types.CoinType{}, fmt.Errorf("can't find any information about coin type: %s", coinName)
		}

		if err = cliCtx.Codec.UnmarshalBinary(bz, &coinType); err != nil {
			return coinType, err
		}
	}

	return coinType, nil
}

func (cliCtx CLIContext) NetInfo() (*ctypes.ResultNetInfo, error) {
	client := cliCtx.Client.(*tmclient.HTTP)
	return client.NetInfo()
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

// Certify verifies the consensus proof at given height
func (cliCtx CLIContext) Certify(height int64) (lite.Commit, error) {
	check, err := tmliteProxy.GetCertifiedCommit(height, cliCtx.Client, cliCtx.Certifier)
	if tmliteErr.IsCommitNotFoundErr(err) {
		return lite.Commit{}, ErrVerifyCommit(height)
	} else if err != nil {
		return lite.Commit{}, err
	}
	return check, nil
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

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType can be app or store
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
