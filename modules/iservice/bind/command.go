package bind

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
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

func SvcBindTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "bind a new service",
		RunE:  cmdr.CreateServiceBindingTx,
	}
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

func (co commander) CreateServiceBindingTx(cmd *cobra.Command, args []string) error {
	//msg := buildSvcDefMsg()
	//ctx := context.NewCoreContextFromViper()
	//res, err := ctx.SignBuildBroadcast(ctx.FromAddressName, msg, co.cdc)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}
