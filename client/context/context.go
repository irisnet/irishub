package context

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/codec"
	cskeys "github.com/irisnet/irishub/crypto/keys"
	"github.com/irisnet/irishub/types"
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
	fromAddress   types.AccAddress
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
		Verifier:      createVerifier(),
		DryRun:        viper.GetBool(client.FlagDryRun),
		GenerateOnly:  viper.GetBool(client.FlagGenerateOnly),
		fromAddress:   fromAddress,
		fromName:      fromName,
		Indent:        viper.GetBool(client.FlagIndentResponse),
	}
}

func createVerifier() tmlite.Verifier {
	trustNodeDefined := viper.IsSet(client.FlagTrustNode)
	if !trustNodeDefined {
		return nil
	}

	trustNode := viper.GetBool(client.FlagTrustNode)
	if trustNode {
		return nil
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

func fromFields(from string) (fromAddr types.AccAddress, fromName string) {
	// In generate-only mode, if the signer key doesn't exist in keystore, the fromAddress can be specified by --from-addr
	if from == "" {
		fromAddrString := viper.GetString(client.FlagFromAddr)
		if fromAddrString == "" {
			return nil, ""
		}
		address, err := types.AccAddressFromBech32(fromAddrString)
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
	if addr, err := types.AccAddressFromBech32(from); err == nil {
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
func (ctx CLIContext) WithCodec(cdc *codec.Codec) CLIContext {
	ctx.Codec = cdc
	return ctx
}

// WithHeight returns a copy of the context with an updated height.
func (ctx CLIContext) WithHeight(height int64) CLIContext {
	ctx.Height = height
	return ctx
}

// WithAccountDecoder returns a copy of the context with an updated account
// decoder.
func (ctx CLIContext) WithAccountDecoder(decoder auth.AccountDecoder) CLIContext {
	ctx.AccDecoder = decoder
	return ctx
}

// WithLogger returns a copy of the context with an updated logger.
func (ctx CLIContext) WithLogger(w io.Writer) CLIContext {
	ctx.Logger = w
	return ctx
}

// WithAccountStore returns a copy of the context with an updated AccountStore.
func (ctx CLIContext) WithAccountStore(accountStore string) CLIContext {
	ctx.AccountStore = accountStore
	return ctx
}

// WithTrustNode returns a copy of the context with an updated TrustNode flag.
func (ctx CLIContext) WithTrustNode(trustNode bool) CLIContext {
	ctx.TrustNode = trustNode
	return ctx
}

// WithNodeURI returns a copy of the context with an updated node URI.
func (ctx CLIContext) WithNodeURI(nodeURI string) CLIContext {
	ctx.NodeURI = nodeURI
	ctx.Client = rpcclient.NewHTTP(nodeURI, "/websocket")
	return ctx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (ctx CLIContext) WithClient(client rpcclient.Client) CLIContext {
	ctx.Client = client
	return ctx
}

// WithUseLedger returns a copy of the context with an updated UseLedger flag.
func (ctx CLIContext) WithUseLedger(useLedger bool) CLIContext {
	ctx.UseLedger = useLedger
	return ctx
}

// WithCertifier - return a copy of the context with an updated Certifier
func (ctx CLIContext) WithCertifier(verifier tmlite.Verifier) CLIContext {
	ctx.Verifier = verifier
	return ctx
}

func (ctx CLIContext) ToMainUnit(coins types.Coins) string {
	ss, _ := ctx.ConvertToMainUnit(coins.String())
	return strings.Join(ss, ",")
}
