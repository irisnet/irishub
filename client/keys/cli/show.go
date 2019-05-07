package keys

import (
	"fmt"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	cryptokeys "github.com/irisnet/irishub/crypto/keys"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tmlibs/cli"
)

const (
	// FlagAddress is the flag for the user's address on the command line.
	FlagAddress = "address"
	// FlagPublicKey represents the user's public key on the command line.
	FlagPublicKey = "pubkey"
	// FlagBechPrefix defines a desired Bech32 prefix encoding for a key.
	FlagBechPrefix = "bech"

	flagMultiSigThreshold  = "multisig-threshold"
	defaultMultiSigKeyName = "multi"
)

var _ cryptokeys.Info = (*multiSigKey)(nil)

type multiSigKey struct {
	name string
	key  crypto.PubKey
}

func (m multiSigKey) GetName() string             { return m.name }
func (m multiSigKey) GetType() cryptokeys.KeyType { return cryptokeys.TypeLocal }
func (m multiSigKey) GetPubKey() crypto.PubKey    { return m.key }
func (m multiSigKey) GetAddress() sdk.AccAddress  { return sdk.AccAddress(m.key.Address()) }

func showKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show key info for the given name",
		Long:  `Return public details of one local key.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  runShowCmd,
	}

	cmd.Flags().String(FlagBechPrefix, "acc", "The Bech32 prefix encoding for a key (acc|val|cons)")
	cmd.Flags().Bool(FlagAddress, false, "output the address only (overrides --output)")
	cmd.Flags().Bool(FlagPublicKey, false, "output the public key only (overrides --output)")
	cmd.Flags().Uint(flagMultiSigThreshold, 1, "K out of N required signatures")
	cmd.Flags().Bool(client.FlagIndentResponse, false, "Add indent to JSON response")

	return cmd
}

func runShowCmd(cmd *cobra.Command, args []string) (err error) {
	var info cryptokeys.Info

	if len(args) == 1 {
		info, err = keys.GetKeyInfo(args[0])
		if err != nil {
			return err
		}
	} else {
		pks := make([]crypto.PubKey, len(args))
		for i, keyName := range args {
			info, err := keys.GetKeyInfo(keyName)
			if err != nil {
				return err
			}
			pks[i] = info.GetPubKey()
		}

		multisigThreshold := viper.GetInt(flagMultiSigThreshold)
		err = validateMultisigThreshold(multisigThreshold, len(args))
		if err != nil {
			return err
		}
		multikey := multisig.NewPubKeyMultisigThreshold(multisigThreshold, pks)
		info = multiSigKey{
			name: defaultMultiSigKeyName,
			key:  multikey,
		}
	}

	isShowAddr := viper.GetBool(FlagAddress)
	isShowPubKey := viper.GetBool(FlagPublicKey)
	isOutputSet := cmd.Flag(cli.OutputFlag).Changed

	if isShowAddr && isShowPubKey {
		return fmt.Errorf("cannot use both --address and --pubkey at once")
	}

	if isOutputSet && (isShowAddr || isShowPubKey) {
		return fmt.Errorf("cannot use --output with --address or --pubkey")
	}

	bechKeyOut, err := keys.GetBechKeyOut(viper.GetString(FlagBechPrefix))
	if err != nil {
		return err
	}

	switch {
	case isShowAddr:
		keys.PrintKeyAddress(info, bechKeyOut)
	case isShowPubKey:
		keys.PrintPubKey(info, bechKeyOut)
	default:
		keys.PrintKeyInfo(info, bechKeyOut)
	}

	return nil
}

func validateMultisigThreshold(k, nKeys int) error {
	if k <= 0 {
		return fmt.Errorf("threshold must be a positive integer")
	}
	if nKeys < k {
		return fmt.Errorf(
			"threshold k of n multisignature: %d < %d", nKeys, k)
	}
	return nil
}
