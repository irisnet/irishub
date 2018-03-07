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
	FlagName        = "svc-name"
	FlagDescription = "svc-description"
)

// nolint
var (
	CmdDefineService = &cobra.Command{
		Use:   "define-service",
		Short: "define a new service",
		RunE:  cmdDefineService,
	}
)

func init() {

	fsService := flag.NewFlagSet("", flag.ContinueOnError)
	fsService.String(FlagName, "", "service name")
	fsService.String(FlagDescription, "", "service description")

	CmdDefineService.Flags().AddFlagSet(fsService)
}

func cmdDefineService(cmd *cobra.Command, args []string) error {
	name := viper.GetString(FlagName)
	desc := viper.GetString(FlagDescription)

	tx := iservice.NewTxDefineService(name, desc)
	return txcmd.DoTx(tx)
}
