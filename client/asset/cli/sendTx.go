package cli

import (
	"encoding/hex"
	"errors"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const DECIMAL_MAX = 18

//GetCmdIssueAsset will transfer the asset to new owner by operators
func GetCmdIssueAsset(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "Issue",
		Example: "iriscli asset issue iris --family=00 --name=IRISnet --symbol=iris --source=00 --initial-supply=100000 --max-supply=1000000000000 --decimal=0 --mintalbe=false  --operators=<account address A>,<account address B>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			//fromStr of the issuer
			fromStr, err := cliCtx.GetFromAddress()
			_=fromStr
			if err != nil {
				return err
			}

			//get the familyStr of asset
			familyStr:=viper.GetString(FlagFamily)
			family,err:= hex.DecodeString(familyStr)
			_=family
			if err != nil {
				return err
			}

			//get the nameStr of asset
			nameStr:= viper.GetString(FlagName)
			if len(nameStr)>32 {
				err := errors.New("name set failed: the length of asset name  must be shorter than 33")
				return err
			}
			reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
			if reg.Match([]byte(nameStr)) {
				err := errors.New("name set failed: the asset name could only include number, english character, and _")
				return err
			}


			//get the symbolStr of asset
			symbolStr:= viper.GetString(FlagSymbol)
			if len(symbolStr)>6 ||  len(symbolStr)<3 {
				err := errors.New("symbol set failed: the length of asset symbol should be shorter than 7 and longer than 2")
				return err
			}
			reg = regexp.MustCompile(`[^a-zA-Z0-9]`)
			if reg.Match([]byte(symbolStr)) {
				err := errors.New("symbol set failed: the asset symbol could only include number and english character")
				return err
			}

			//get the sourceStr of asset
			sourceStr:= viper.GetString(FlagSource)
			_=sourceStr

			//get the initialSupply of asset
			initialSupplyInt64 := viper.GetInt64(FlagInitialSupply)
			if initialSupplyInt64 < 0 {
				err := errors.New("initial-supply set failed: the initial supply of asset must be positive")
				return err
			}
			if initialSupplyInt64 > 1000000000000 {
				err := errors.New("initial-supply set failed: the initial supply of asset must be smaller than 100 billion")
				return err
			}
			//initialSupply of asset
			initialSupply := uint64(initialSupplyInt64)
			_ = initialSupply

			//get the maxSupply of asset
			maxSupplyInt64 := viper.GetInt64(FlagMaxSupply)
			if maxSupplyInt64 < 0 {
				err := errors.New("max-supply set failed: the max supply of asset must be positive")
				return err
			}
			//maxSupply of asset
			maxSupply:= uint64(maxSupplyInt64)
			_=maxSupply

			//get the mintable of asset
			mintable:=viper.GetBool(FlagMintable)
			_=mintable

			//get the decimal of asset
			DecimalInt64 := viper.GetInt64(FlagDecimal)
			if DecimalInt64 < 0 {
				err := errors.New("decimal set failed: the decimal  of asset must be positive")
				return err
			}
			if DecimalInt64 > 18 {
				err := errors.New("decimal set failed: the token can have a maximum of 18 digits of decimal")
				return err
			}
			//decimal
			_ = uint8(DecimalInt64)

			//get the operator of the asset
			operatorsStrs := viper.GetStringSlice(FlagOperators)
			var operators []sdk.AccAddress

			for _, operatorsStr := range operatorsStrs {

				ope, err := sdk.AccAddressFromBech32(operatorsStr)
				if err != nil {
					return err
				}
				operators=append(operators, ope)

			}

			//asset msg has not been complished
			var msg sdk.Msg
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(FlagFamily, "00", "The family of the newly issued asset,00:fungible,01:non-fungible")
	cmd.Flags().String(FlagName, "", "The name of the newly issued asset.The length of asset name  must be shorter than 33. The asset name could only include number, english character, and _")
	cmd.Flags().String(FlagSymbol, "", "The symbol representing the newly issued asset. The length of asset symbol should be shorter than 7 and longer than 2. The asset symbol could only include number and english character")
	cmd.Flags().String(FlagSource, "00", "00:native,01:external,gatewayID")
	cmd.Flags().String(FlagInitialSupply, "", "The initial supply for this token")
	cmd.Flags().String(FlagMaxSupply, "1000,000,000,000", "The hard crap of this asset")
	cmd.Flags().String(FlagDecimal, "0", "The decimal of this asset")
	cmd.Flags().String(FlagMintable, "false", "Whether this token could be minted")
	cmd.Flags().String(FlagOperators, "", "The operator set of this asset")
	cmd.MarkFlagRequired(FlagName)
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagInitialSupply)
	return cmd
}
