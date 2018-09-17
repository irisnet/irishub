package keys

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ccrypto "github.com/cosmos/cosmos-sdk/crypto"
	cryptokeys "github.com/cosmos/cosmos-sdk/crypto/keys"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagType     = "type"
	flagRecover  = "recover"
	flagNoBackup = "no-backup"
	flagDryRun   = "dry-run"
	flagAccount  = "account"
	flagIndex    = "index"
)

func addKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a new key, or import from seed",
		Long: `Add a public/private key pair to the key store.
If you select --seed/-s you can recover a key from the seed
phrase, otherwise, a new key will be generated.`,
		RunE: runAddCmd,
	}
	cmd.Flags().StringP(flagType, "t", "secp256k1", "Type of private key (secp256k1|ed25519)")
	cmd.Flags().Bool(client.FlagUseLedger, false, "Store a local reference to a private key on a Ledger device")
	cmd.Flags().Bool(flagRecover, false, "Provide seed phrase to recover existing key instead of creating")
	cmd.Flags().Bool(flagNoBackup, false, "Don't print out seed phrase (if others are watching the terminal)")
	cmd.Flags().Bool(flagDryRun, false, "Perform action, but don't add key to local keystore")
	cmd.Flags().Uint32(flagAccount, 0, "Account number for HD derivation")
	cmd.Flags().Uint32(flagIndex, 0, "Index number for HD derivation")
	return cmd
}

// nolint: gocyclo
// TODO remove the above when addressing #1446
func runAddCmd(cmd *cobra.Command, args []string) error {
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
		kb, err = keys.GetKeyBase()
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

		// ask for a password when generating a local key
		if !viper.GetBool(client.FlagUseLedger) {
			pass, err = keys.GetCheckPassword(
				"Enter a passphrase for your key:",
				"Repeat the passphrase:", buf)
			if err != nil {
				return err
			}
		}
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
		keys.PrintInfo(cdc, info)
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
		json, err := cdc.MarshalJSON(out)
		if err != nil {
			panic(err) // really shouldn't happen...
		}
		fmt.Println(string(json))
	default:
		panic(fmt.Sprintf("I can't speak: %s", output))
	}
}
