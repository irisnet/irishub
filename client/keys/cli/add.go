package keys

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	ccrypto "github.com/irisnet/irishub/crypto"
	cryptokeys "github.com/irisnet/irishub/crypto/keys"
	"github.com/irisnet/irishub/crypto/keystore"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagType     = "type"
	flagRecover  = "recover"
	flagNoBackup = "no-backup"
	flagDryRun   = "dry-run"
	flagAccount  = "account"
	flagIndex    = "index"
	flagMultisig = "multisig"
	flagNoSort   = "nosort"
	flagKeystore = "keystore"
)

func addKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a new key, or import from seed",
		Long: `Add a public/private key pair to the key store.
If you select --seed/-s you can recover a key from the seed
phrase, otherwise, a new key will be generated.`,
		Example: "iriscli keys add <key name>",
		RunE:    runAddCmd,
	}
	cmd.Flags().StringSlice(flagMultisig, nil, "Construct and store a multisig public key (implies --pubkey)")
	cmd.Flags().Uint(flagMultiSigThreshold, 1, "K out of N required signatures. For use in conjunction with --multisig")
	cmd.Flags().Bool(flagNoSort, false, "Keys passed to --multisig are taken in the order they're supplied")
	cmd.Flags().String(FlagPublicKey, "", "Parse a public key in bech32 format and save it to disk")
	cmd.Flags().StringP(flagType, "t", "secp256k1", "Type of private key (secp256k1|ed25519)")
	cmd.Flags().Bool(client.FlagUseLedger, false, "Store a local reference to a private key on a Ledger device")
	cmd.Flags().Bool(flagRecover, false, "Provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flagKeystore, "", "Provide keystore file to recover existing key instead of creating. For use in conjunction with --recover")
	cmd.Flags().Bool(flagNoBackup, false, "Don't print out seed phrase (if others are watching the terminal)")
	cmd.Flags().Bool(flagDryRun, false, "Perform action, but don't add key to local keystore")
	cmd.Flags().Uint32(flagAccount, 0, "Account number for HD derivation")
	cmd.Flags().Uint32(flagIndex, 0, "Index number for HD derivation")
	cmd.Flags().Bool(client.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

// nolint: gocyclo
// TODO remove the above when addressing #1446
func runAddCmd(_ *cobra.Command, args []string) error {
	var kb cryptokeys.Keybase
	var err error
	var name, pass string

	buf := keys.BufferStdin()
	if viper.GetBool(flagDryRun) {
		// we throw this away, so don't enforce args,
		// we want to get a new random seed phrase quickly
		kb = keys.MockKeyBase()
		pass = "throwing-this-key-away"
		name = "inmemorykey"
	} else {
		if len(args) != 1 || len(args[0]) == 0 {
			return errors.New("you must provide a name for the key")
		}
		name = args[0]
		kb, err = keys.GetKeyBaseWithWritePerm()
		if err != nil {
			return err
		}

		_, err := kb.Get(name)
		if err == nil {
			// account exists, ask for user confirmation
			if response, err := keys.GetConfirmation(
				fmt.Sprintf("override the existing name %s", name), buf); err != nil || !response {
				return err
			}
		}

		multisigKeys := viper.GetStringSlice(flagMultisig)
		if len(multisigKeys) != 0 {
			var pks []crypto.PubKey

			multisigThreshold := viper.GetInt(flagMultiSigThreshold)
			if err := validateMultisigThreshold(multisigThreshold, len(multisigKeys)); err != nil {
				return err
			}

			addressMap := make(map[string]interface{})
			for _, keyname := range multisigKeys {
				k, err := kb.Get(keyname)
				if err != nil {
					return err
				}
				if k.GetType() == cryptokeys.TypeMulti {
					return errors.New("can not create multi-sign key with other multi-sign keys")
				}
				if _, ok := addressMap[k.GetAddress().String()]; ok {
					return fmt.Errorf("can not create multi-sign key with duplicate address %s", k.GetAddress().String())
				}
				addressMap[k.GetAddress().String()] = nil
				pks = append(pks, k.GetPubKey())
			}

			// Handle --nosort
			if !viper.GetBool(flagNoSort) {
				sort.Slice(pks, func(i, j int) bool {
					return bytes.Compare(pks[i].Address(), pks[j].Address()) < 0
				})
			}

			pk := multisig.NewPubKeyMultisigThreshold(multisigThreshold, pks)
			if _, err := kb.CreateMulti(name, pk); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Key %q saved to disk.\n", name)
			return nil
		}

		// ask for a password when generating a local key
		if viper.GetString(FlagPublicKey) == "" && !viper.GetBool(client.FlagUseLedger) {
			pass, err = keys.GetCheckPassword(
				"Enter a passphrase for your key:",
				"Repeat the passphrase:", buf)
			if err != nil {
				return err
			}
		}
	}

	if viper.GetString(FlagPublicKey) != "" {
		pk, err := sdk.GetAccPubKeyBech32(viper.GetString(FlagPublicKey))
		if err != nil {
			return err
		}
		_, err = kb.CreateOffline(name, pk)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Key %q saved to disk.\n", name)
		return nil
	}

	if viper.GetBool(client.FlagUseLedger) {
		account := uint32(viper.GetInt(flagAccount))
		index := uint32(viper.GetInt(flagIndex))
		path := ccrypto.DerivationPath{44, 118, account, 0, index}
		algo := cryptokeys.SigningAlgo(viper.GetString(flagType))
		info, err := kb.CreateLedger(name, path, algo)
		if err != nil {
			return err
		}
		printCreate(info, "")
	} else if viper.GetBool(flagRecover) {
		keystoreFile := viper.GetString(flagKeystore)
		if len(keystoreFile) > 0 {
			var passphrase string
			if !keys.InputIsTty() {
				passphrase = pass
			} else {
				prompt := fmt.Sprintf("Password of the keystore file:")
				passphrase, err = keys.GetPassword(prompt, buf)
				if err != nil {
					return fmt.Errorf("Error reading passphrase: %v", err)
				}
			}

			km, err := keystore.NewKeyStoreKeyManager(keystoreFile, passphrase)
			if err != nil {
				return err
			}
			info, err := kb.ImportPrivateKey(name, pass, km.GetPrivKey())
			if err != nil {
				return err
			}
			// print out results without the seed phrase
			viper.Set(flagNoBackup, true)
			printCreate(info, "")
		} else {
			seed, err := keys.GetSeed(
				"Enter your recovery seed phrase:", buf)
			if err != nil {
				return err
			}
			info, err := kb.CreateKey(name, seed, pass)
			if err != nil {
				return err
			}
			// print out results without the seed phrase
			viper.Set(flagNoBackup, true)
			printCreate(info, "")
		}
	} else {
		algo := cryptokeys.SigningAlgo(viper.GetString(flagType))
		info, seed, err := kb.CreateMnemonic(name, cryptokeys.English, pass, algo)
		if err != nil {
			return err
		}
		printCreate(info, seed)
	}
	return nil
}

func printCreate(info cryptokeys.Info, seed string) {
	output := viper.Get(cli.OutputFlag)
	switch output {
	case "text":
		keys.PrintKeyInfo(info, keys.Bech32KeyOutput)
		// print seed unless requested not to.
		if !viper.GetBool(client.FlagUseLedger) && !viper.GetBool(flagNoBackup) {
			fmt.Println("**Important** write this seed phrase in a safe place.")
			fmt.Println("It is the only way to recover your account if you ever forget your password.")
			fmt.Println()
			fmt.Println(seed)
		}
	case "json":
		out, err := keys.Bech32KeyOutput(info)
		if err != nil {
			panic(err)
		}
		if !viper.GetBool(flagNoBackup) {
			out.Seed = seed
		}
		var jsonString []byte
		if viper.GetBool(client.FlagIndentResponse) {
			jsonString, err = cdc.MarshalJSONIndent(out, "", "  ")
		} else {
			jsonString, err = cdc.MarshalJSON(out)
		}
		if err != nil {
			panic(err) // really shouldn't happen...
		}
		fmt.Println(string(jsonString))
	default:
		panic(fmt.Sprintf("I can't speak: %s", output))
	}
}
