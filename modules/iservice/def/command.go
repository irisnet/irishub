package def

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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

func SvcDefTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "define",
		Short: "define a new service",
		RunE:  cmdr.createServiceDefinitionTx,
	}
	//application args
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

	cmd.Flags().AddFlagSet(fsService)
	return cmd
}

type commander struct {
	cdc *wire.Codec
}

func (co commander) createServiceDefinitionTx(cmd *cobra.Command, args []string) error {
	msg := buildSvcDefMsg()
	ctx := context.NewCoreContextFromViper()

	//txBytes, _ := ctx.SignAndBuild(ctx.FromAddressName, "1234567890", msg, co.cdc)
	//
	//res, err :=  ctx.BroadcastTx(txBytes)

	res, err := ctx.SignBuildBroadcast(ctx.FromAddressName, msg, co.cdc)
	if err != nil {
		return err
	}

	fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func buildSvcDefMsg() SvcDefMsg {
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

	msg := SvcDefMsg{
		Name:        name,
		Description: desc,
		Tags:        tags,
		Creator:     creator,
		ChainID:     chainID,
		Messaging:   messaging,
		Methods: []Method{
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
	return msg
}
