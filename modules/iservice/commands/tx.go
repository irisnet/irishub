package commands

import (
	txcmd "github.com/cosmos/cosmos-sdk/client/commands/txs"
	"github.com/irisnet/iris-hub/modules/iservice"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// nolint
const (
	FlagName                 = "svc-name"
	FlagDescription          = "svc-description"
	FlagTags                 = "svc-tags"
	FlagCreator              = "svc-creator"
	FlagChainID              = "svc-chainID"
	FlagMessaging            = "svc-messaging"
	FlagMethodsID            = "svc-method-id"
	FlagMethodsName          = "svc-method-name"
	FlagMethodsDescription   = "svc-method-desc"
	FlagMethodsInput         = "svc-method-input"
	FlagMethodsOutput        = "svc-method-output"
	FlagMethodsError         = "svc-method-error"
	FlagMethodsOutputPrivacy = "svc-method-outputPrivacy"
)

// nolint
var (
	CmdCreateServiceDefinitionTx = &cobra.Command{
		Use:   "define-service",
		Short: "define a new service",
		RunE:  createServiceDefinitionTx,
	}
)

func init() {

	fsService := flag.NewFlagSet("", flag.ContinueOnError)
	fsService.String(FlagName, "", "service name")
	fsService.String(FlagDescription, "", "service description")
	fsService.String(FlagTags, "", "service tags")
	fsService.String(FlagCreator, "", "service creator")
	fsService.String(FlagChainID, "", "service chainID")
	fsService.String(FlagMessaging, "", "service messaging")
	fsService.Int(FlagMethodsID, 0, "service method-id")
	fsService.String(FlagMethodsName, "", "service method-name")
	fsService.String(FlagMethodsDescription, "", "service method-desc")
	fsService.String(FlagMethodsInput, "", "service method-input")
	fsService.String(FlagMethodsOutput, "", "service method-output")
	fsService.String(FlagMethodsError, "", "service method-error")
	fsService.String(FlagMethodsOutputPrivacy, "", "service method-outputPrivacy")

	CmdCreateServiceDefinitionTx.Flags().AddFlagSet(fsService)
}

func createServiceDefinitionTx(cmd *cobra.Command, args []string) error {
	name := viper.GetString(FlagName)
	desc := viper.GetString(FlagDescription)
	tags := viper.GetString(FlagTags)
	creator := viper.GetString(FlagCreator)
	chainID := viper.GetString(FlagChainID)
	messaging := viper.GetString(FlagMessaging)
	methodsID := viper.GetInt(FlagMethodsID)
	methodsName := viper.GetString(FlagMethodsName)
	methodsDescription := viper.GetString(FlagMethodsDescription)
	methodsInput := viper.GetString(FlagMethodsInput)
	methodsOutput := viper.GetString(FlagMethodsOutput)
	methodsError := viper.GetString(FlagMethodsError)
	methodsOutputPrivacy := viper.GetString(FlagMethodsOutputPrivacy)

	service := iservice.ServiceDefinition{
		Name: name,
		Description: desc,
		Tags: tags,
		Creator: creator,
		ChainID: chainID,
		Messaging: messaging,
		Methods: []iservice.ServiceMethod{
			{
				ID:            methodsID,
				Name:          methodsName,
				Description:   methodsDescription,
				Input:         methodsInput,
				Output:        methodsOutput,
				Error:         methodsError,
				OutputPrivacy: methodsOutputPrivacy,
			},
		},
	}

	tx := iservice.NewTxDefineService(service)
	return txcmd.DoTx(tx)
}
