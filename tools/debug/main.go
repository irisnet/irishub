package debug

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	iris "github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func init() {
	RootCmd.AddCommand(txCmd)
	RootCmd.AddCommand(pubkeyCmd)
	RootCmd.AddCommand(addrCmd)
	RootCmd.AddCommand(rawBytesCmd)
	RootCmd.AddCommand(randSecretCmd)
	RootCmd.AddCommand(hashLockCmd)
}

var RootCmd = &cobra.Command{
	Use:          "debug",
	Short:        "Iris debug tool",
	SilenceUsage: true,
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Decode a iris tx from hex or base64",
	RunE:  runTxCmd,
}

var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Decode a pubkey from hex, base64, or bech32",
	RunE:  runPubKeyCmd,
}

var addrCmd = &cobra.Command{
	Use:   "addr",
	Short: "Convert an address between hex and bech32",
	RunE:  runAddrCmd,
}

var rawBytesCmd = &cobra.Command{
	Use:   "raw-bytes",
	Short: "Convert raw bytes output (eg. [10 21 13 255]) to hex",
	RunE:  runRawBytesCmd,
}

var randSecretCmd = &cobra.Command{
	Use:   "rand-secret",
	Short: "Generate a random secret",
	RunE:  runRandSecretCmd,
}

var hashLockCmd = &cobra.Command{
	Use:   "hash-lock",
	Short: "Generate a hash lock with secret and timestamp(if provided)",
	RunE:  runHashLockCmd,
}

func runRandSecretCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("Expected no arg")
	}
	b := make([]byte, 32)
	m, err := rand.Read(b)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", hex.EncodeToString(b[:m]))
	return nil
}

func runHashLockCmd(cmd *cobra.Command, args []string) error {
	lenArgs := len(args)
	if lenArgs != 1 && lenArgs != 2 {
		return fmt.Errorf("Expected one or two args")
	}
	if len(args[0]) != 64 {
		return fmt.Errorf("Expected bech32")
	}
	secret, err := hex.DecodeString(args[0])
	if err != nil {
		return err
	}

	hashLock := []byte{}
	switch lenArgs {
	case 1:
		hashLock = sdk.SHA256(secret)
	case 2:
		timestamp, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
		hashLock = sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}
	hashLockHexStr := hex.EncodeToString(hashLock)
	fmt.Printf("%s\n", hashLockHexStr)
	return nil
}

func runRawBytesCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}
	stringBytes := args[0]
	stringBytes = strings.Trim(stringBytes, "[")
	stringBytes = strings.Trim(stringBytes, "]")
	spl := strings.Split(stringBytes, " ")

	byteArray := []byte{}
	for _, s := range spl {
		b, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		byteArray = append(byteArray, byte(b))
	}
	fmt.Printf("%X\n", byteArray)
	return nil
}

func runPubKeyCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	pubkeyString := args[0]
	var pubKeyI crypto.PubKey

	// try hex, then base64, then bech32
	pubkeyBytes, err := hex.DecodeString(pubkeyString)
	if err != nil {
		var err2 error
		pubkeyBytes, err2 = base64.StdEncoding.DecodeString(pubkeyString)
		if err2 != nil {
			var err3 error
			pubKeyI, err3 = sdk.GetAccPubKeyBech32(pubkeyString)
			if err3 != nil {
				var err4 error
				pubKeyI, err4 = sdk.GetValPubKeyBech32(pubkeyString)

				if err4 != nil {
					var err5 error
					pubKeyI, err5 = sdk.GetConsPubKeyBech32(pubkeyString)
					if err5 != nil {
						return fmt.Errorf(`Expected hex, base64, or bech32. Got errors:
								hex: %v,
								base64: %v
								bech32 Acc: %v
								bech32 Val: %v
								bech32 Cons: %v`,
							err, err2, err3, err4, err5)
					}

				}
			}

		}
	}

	var pubKey ed25519.PubKeyEd25519
	if pubKeyI == nil {
		copy(pubKey[:], pubkeyBytes)
	} else {
		pubKey = pubKeyI.(ed25519.PubKeyEd25519)
		pubkeyBytes = pubKey[:]
	}

	cdc := iris.MakeLatestCodec()
	pubKeyJSONBytes, err := cdc.MarshalJSON(pubKey)
	if err != nil {
		return err
	}
	accPub, err := sdk.Bech32ifyAccPub(pubKey)
	if err != nil {
		return err
	}
	valPub, err := sdk.Bech32ifyValPub(pubKey)
	if err != nil {
		return err
	}
	consPub, err := sdk.Bech32ifyConsPub(pubKey)
	if err != nil {
		return err
	}
	fmt.Println("Address:", pubKey.Address())
	fmt.Printf("Hex: %X\n", pubkeyBytes)
	fmt.Println("JSON (base64):", string(pubKeyJSONBytes))
	fmt.Println("Bech32 Acc:", accPub)
	fmt.Println("Bech32 Val:", valPub)
	fmt.Println("Bech32 Cons:", consPub)
	return nil
}

func runAddrCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	addrString := args[0]
	var addr []byte

	// try hex, then bech32
	var err error
	addr, err = hex.DecodeString(addrString)
	if err != nil {
		var err2 error
		addr, err2 = sdk.AccAddressFromBech32(addrString)
		if err2 != nil {
			var err3 error
			addr, err3 = sdk.ValAddressFromBech32(addrString)

			if err3 != nil {
				var err4 error
				addr, err4 = sdk.ConsAddressFromBech32(addrString)
				if err4 != nil {
					return fmt.Errorf(`Expected hex or bech32. Got errors:
							hex: %v,
							bech32 Acc: %v
							bech32 Val: %v
							bech32 Cons: %v:`,
						err, err2, err3, err4)
				}

			}
		}
	}

	accAddr := sdk.AccAddress(addr)
	valAddr := sdk.ValAddress(addr)
	consAddr := sdk.ConsAddress(addr)

	fmt.Printf("Address (Hex): %X\n", addr)
	fmt.Printf("Bech32 Acc: %s\n", accAddr)
	fmt.Printf("Bech32 Val: %s\n", valAddr)
	fmt.Printf("Bech32 Cons: %s\n", consAddr)
	return nil
}

func runTxCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	txString := args[0]

	// try hex, then base64
	txBytes, err := hex.DecodeString(txString)
	if err != nil {
		var err2 error
		txBytes, err2 = base64.StdEncoding.DecodeString(txString)
		if err2 != nil {
			return fmt.Errorf(`Expected hex or base64. Got errors:
			hex: %v,
			base64: %v
			`, err, err2)
		}
	}

	var tx = auth.StdTx{}
	cdc := iris.MakeLatestCodec()

	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return err
	}

	bz, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, bz, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}
