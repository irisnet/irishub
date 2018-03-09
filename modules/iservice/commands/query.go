package commands

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/search"
	"sort"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"errors"
	is "github.com/irisnet/iris-hub/modules/iservice"
	"strings"
)

//nolint
var (
	CmdQueryServiceDefinition = &cobra.Command{
		Use:   "service-definition",
		RunE:  cmdQueryServiceDefinition,
		Short: "Query a service definition based on name",
	}

	FlagServiceName   = "svc-name"
	FlagChainId       = "chain-id"
	FlagMessagingType = "msg-type"
	FlagAddress       = "address"
	FlagSvTags        = "svc-tags"
)

func init() {
	//Add Flags
	fsSn := flag.NewFlagSet("", flag.ContinueOnError)
	fsSn.String(FlagServiceName, "", "Name of the service")
	fsSn.String(FlagChainId, "", "ChainId of the service")
	fsSn.String(FlagMessagingType, "", "MessagingType of the service")
	fsSn.String(FlagAddress, "", "publisher's address of the service")
	fsSn.String(FlagSvTags, "", "tags of the service")

	CmdQueryServiceDefinition.Flags().AddFlagSet(fsSn)
}

func cmdQueryServiceDefinition(cmd *cobra.Command, args []string) error {

	prove := !viper.GetBool(commands.FlagTrustNode)
	queries := []string{}

	name := viper.GetString(FlagServiceName)
	if name != "" {
		queries = append(queries, is.QServiceName(name))
	}

	chainId := viper.GetString(FlagChainId)
	if chainId != "" {
		queries = append(queries, is.QChainId(chainId))
	}

	messagingType := viper.GetString(FlagMessagingType)
	if messagingType != "" {
		queries = append(queries, is.QMessagingType(messagingType))
	}

	address := viper.GetString(FlagAddress)
	if address != "" {
		queries = append(queries, is.QSender(address))
	}

	tags := viper.GetString(FlagSvTags)
	tx_tags := strings.Split(tags, ",")
	for _, tag := range tx_tags {
		kv := strings.Split(tag, "=")
		if len(kv) == 2 {
			queries = append(queries, is.QKVTag(kv[0], kv[1]))
		} else {
			queries = append(queries, is.QKeyTag(kv[0]))
		}
	}

	all, err := findAnyTx(prove, queries)
	output, err := search.FormatSearch(all, is.ExtractSvDefineTx)
	if err != nil {
		return errors.New("record is empty")
	}

	return search.Output(output)
}

func findAnyTx(prove bool, queries []string) ([]*ctypes.ResultTx, error) {
	var all []*ctypes.ResultTx
	// combine all requests
	for _, q := range queries {
		txs, err := search.FindTx(q, prove)
		if err != nil {
			return nil, err
		}
		all = append(all, txs...)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Height > all[j].Height
	})
	return all, nil
}
