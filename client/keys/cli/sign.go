package keys

import (
	"encoding/base64"
	"fmt"
	"github.com/irisnet/irishub/client/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagFrom = "from"
	flagTx   = "tx"
)

func init() {
	keySignCmd.Flags().String(flagFrom, "", "Name of private key with which to sign")
	keySignCmd.Flags().String(flagTx, "", "Tx data for sign")
}

var keySignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign user specified data",
	Long:  `Sign user specified data and return a base64 encoding signature`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := viper.GetString(flagFrom)
		tx := viper.GetString(flagTx)

		kb, err := keys.GetKeyBase()
		if err != nil {
			return err
		}

		buf := keys.BufferStdin()
		password, err := keys.GetPassword(
			"Enter passphrase:", buf)
		if err != nil {
			return err
		}

		sig, _, err := kb.Sign(name, password, []byte(tx))
		if err != nil {
			return err
		}
		encoded := base64.StdEncoding.EncodeToString(sig)
		fmt.Println(string(encoded))
		return nil
	},
}
