package commands

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	"github.com/cosmos/cosmos-sdk/stack"
	"github.com/irisnet/iris-hub/modules/iservice"
)

//nolint
var (
	CmdQueryServiceDefinition = &cobra.Command{
		Use:   "service-definition",
		RunE:  cmdQueryServiceDefinition,
		Short: "Query a service definition based on name",
	}

	FlagServiceName = "svc-name"
)

func init() {
	//Add Flags
	fsSn := flag.NewFlagSet("", flag.ContinueOnError)
	fsSn.String(FlagServiceName, "", "Name of the service")

	CmdQueryServiceDefinition.Flags().AddFlagSet(fsSn)
}

func cmdQueryServiceDefinition(cmd *cobra.Command, args []string) error {

	name := viper.GetString(FlagServiceName)
	prove := !viper.GetBool(commands.FlagTrustNode)
	key := stack.PrefixedKey(iservice.Name(), iservice.GetServiceDefinitionKey(name))
	var svc iservice.ServiceDefinition
	height, err := query.GetParsed(key, &svc, query.GetHeight(), prove)
	if err != nil {
		return err
	}

	return query.OutputProof(svc, height)
}
