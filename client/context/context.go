package context

import (
	"bytes"
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap"
	"io"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/codec"
	cskeys "github.com/irisnet/irishub/crypto/keys"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmlite "github.com/tendermint/tendermint/lite"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// CLIContext implements a typical CLI context created in SDK modules for
// transaction handling and queries.
type CLIContext struct {
	Codec         *codec.Codec
	AccDecoder    auth.AccountDecoder
	Client        rpcclient.Client
	Logger        io.Writer
	OutputFormat  string
	Height        int64
	NodeURI       string
	AccountStore  string
	TrustNode     bool
	UseLedger     bool
	Commit        bool
	Async         bool
	JSON          bool
	PrintResponse bool
	Verifier      tmlite.Verifier
	GenerateOnly  bool
	fromAddress   sdk.AccAddress
	fromName      string
	Indent        bool
	DryRun        bool
}

// NewCLIContext returns a new initialized CLIContext with parameters from the
// command line using Viper.
func NewCLIContext() CLIContext {
	var rpc rpcclient.Client

	from := viper.GetString(client.FlagFrom)
	fromAddress, fromName := fromFields(from)

	nodeURI := viper.GetString(client.FlagNode)
	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}

	return CLIContext{
		Client:        rpc,
		NodeURI:       nodeURI,
		AccountStore:  protocol.AccountStore,
		Height:        viper.GetInt64(client.FlagHeight),
		OutputFormat:  viper.GetString(cli.OutputFlag),
		TrustNode:     viper.GetBool(client.FlagTrustNode),
		UseLedger:     viper.GetBool(client.FlagUseLedger),
		Async:         viper.GetBool(client.FlagAsync),
		Commit:        viper.GetBool(client.FlagCommit),
		JSON:          viper.GetBool(client.FlagJson),
		PrintResponse: viper.GetBool(client.FlagPrintResponse),
		Verifier:      createVerifier(rpc),
		DryRun:        viper.GetBool(client.FlagDryRun),
		GenerateOnly:  viper.GetBool(client.FlagGenerateOnly),
		fromAddress:   fromAddress,
		fromName:      fromName,
		Indent:        viper.GetBool(client.FlagIndentResponse),
	}
}

func createVerifier(rpc rpcclient.SignClient) tmlite.Verifier {
	trustNodeDefined := viper.IsSet(client.FlagTrustNode)
	if !trustNodeDefined {
		return nil
	}

	trustNode := viper.GetBool(client.FlagTrustNode)
	if trustNode {
		return nil
	} else {
		height := int64(1)
		if _, err := rpc.Commit(&height); err != nil {
			fmt.Printf("snapshot's node can't verify the proof of result, you must set '--trust-node=true'\n")
			os.Exit(1)
		}
	}

	chainID := viper.GetString(client.FlagChainID)
	home := viper.GetString(cli.HomeFlag)
	nodeURI := viper.GetString(client.FlagNode)

	var errMsg bytes.Buffer
	if chainID == "" {
		errMsg.WriteString("--chain-id ")
	}
	if home == "" {
		errMsg.WriteString("--home ")
	}
	if nodeURI == "" {
		errMsg.WriteString("--node ")
	}
	if errMsg.Len() != 0 {
		fmt.Printf("Must specify these options: %s when --trust-node is false\n", errMsg.String())
		os.Exit(1)
	}

	node := rpcclient.NewHTTP(nodeURI, "/websocket")
	cacheSize := 10 // TODO: determine appropriate cache size
	verifier, err := tmliteProxy.NewVerifier(
		chainID, home,
		node, log.NewNopLogger(), cacheSize,
	)

	if err != nil {
		fmt.Printf("Create verifier failed: %s\n", err.Error())
		fmt.Printf("Please check network connection and verify the address of the node to connect to\n")
		os.Exit(1)
	}

	return verifier
}

func fromFields(from string) (fromAddr sdk.AccAddress, fromName string) {
	// In generate-only mode, if the signer key doesn't exist in keystore, the fromAddress can be specified by --from-addr
	if from == "" {
		fromAddrString := viper.GetString(client.FlagFromAddr)
		if fromAddrString == "" {
			return nil, ""
		}
		address, err := sdk.AccAddressFromBech32(fromAddrString)
		if err != nil {
			fmt.Printf("invalid from address %s\n", fromAddrString)
			os.Exit(1)
		}
		fromAddr = address
		fromName = ""
		return
	}

	keybase, err := keys.GetKeyBase()
	if err != nil {
		fmt.Println("no keybase found")
		os.Exit(1)
	}

	var info cskeys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			fmt.Printf("could not find key %s\n", from)
			os.Exit(1)
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("could not find key %s\n", from))
			os.Exit(1)
		}
	}

	fromAddr = info.GetAddress()
	fromName = info.GetName()
	return
}

// WithCodec returns a copy of the context with an updated codec.
func (cliCtx CLIContext) WithCodec(cdc *codec.Codec) CLIContext {
	cliCtx.Codec = cdc
	return cliCtx
}

// WithHeight returns a copy of the context with an updated height.
func (cliCtx CLIContext) WithHeight(height int64) CLIContext {
	cliCtx.Height = height
	return cliCtx
}

// WithAccountDecoder returns a copy of the context with an updated account
// decoder.
func (cliCtx CLIContext) WithAccountDecoder(decoder auth.AccountDecoder) CLIContext {
	cliCtx.AccDecoder = decoder
	return cliCtx
}

// WithLogger returns a copy of the context with an updated logger.
func (cliCtx CLIContext) WithLogger(w io.Writer) CLIContext {
	cliCtx.Logger = w
	return cliCtx
}

// WithAccountStore returns a copy of the context with an updated AccountStore.
func (cliCtx CLIContext) WithAccountStore(accountStore string) CLIContext {
	cliCtx.AccountStore = accountStore
	return cliCtx
}

// WithTrustNode returns a copy of the context with an updated TrustNode flag.
func (cliCtx CLIContext) WithTrustNode(trustNode bool) CLIContext {
	cliCtx.TrustNode = trustNode
	return cliCtx
}

// WithNodeURI returns a copy of the context with an updated node URI.
func (cliCtx CLIContext) WithNodeURI(nodeURI string) CLIContext {
	cliCtx.NodeURI = nodeURI
	cliCtx.Client = rpcclient.NewHTTP(nodeURI, "/websocket")
	return cliCtx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (cliCtx CLIContext) WithClient(client rpcclient.Client) CLIContext {
	cliCtx.Client = client
	return cliCtx
}

// WithUseLedger returns a copy of the context with an updated UseLedger flag.
func (cliCtx CLIContext) WithUseLedger(useLedger bool) CLIContext {
	cliCtx.UseLedger = useLedger
	return cliCtx
}

// WithCertifier - return a copy of the context with an updated Certifier
func (cliCtx CLIContext) WithCertifier(verifier tmlite.Verifier) CLIContext {
	cliCtx.Verifier = verifier
	return cliCtx
}

func (cliCtx CLIContext) GetCoinType(coinName string) (sdk.CoinType, error) {
	var coinType sdk.CoinType
	coinName = strings.ToLower(coinName)
	if coinName == "" {
		return sdk.CoinType{}, fmt.Errorf("coin name is empty")
	}
	if coinName == sdk.Iris {
		coinType = sdk.IrisCoinType
	} else if strings.HasPrefix(coinName, coinswap.FormatUniABSPrefix) {
		return coinswap.GetUniCoinType(coinName)
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
	if len(coinsStr) == 0 {
		return coins, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		mainUnit, err := sdk.GetCoinName(coinStr)
		coinType, err := cliCtx.GetCoinType(mainUnit)
		if err != nil {
			coins = append(coins, coinStr)
			continue
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
	if err != nil {
		return sdk.Coin{}, err
	}
	coinType, err := cliCtx.GetCoinType(mainUnit)
	if err != nil {
		return sdk.ParseCoin(coinStr)
	}

	coin, err := coinType.ConvertToMinDenomCoin(coinStr)
	if err != nil {
		return sdk.Coin{}, err
	}
	return coin, nil
}

func (cliCtx CLIContext) ParseCoins(coinsStr string) (coins sdk.Coins, err error) {
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
			coinMap[coin.Denom] = coinMap[coin.Denom].Add(coin)
		} else {
			coinMap[coin.Denom] = coin
		}
	}

	for _, coin := range coinMap {
		coins = append(coins, coin)
	}

	return sdk.NewCoins(coins...), nil
}

func (cliCtx CLIContext) ToMainUnit(coins sdk.Coins) string {
	ss, _ := cliCtx.ConvertToMainUnit(coins.String())
	return strings.Join(ss, ",")
}
