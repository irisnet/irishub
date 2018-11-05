package client

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

type cliConfig struct {
	Home      string `toml:"home"`
	ChainID   string `toml:"chain-id"`
	TrustNode bool   `toml:"trust-node"`
	Encoding  string `toml:"encoding"`
	Output    string `toml:"output"`
	Node      string `toml:"node"`
	Trace     bool   `toml:"trace"`
	Fee       string `toml:"fee"`
}

// ConfigCmd returns a CLI command to interactively create a
// Gaia CLI config file.
func ConfigCmd() *cobra.Command {
	cfg := &cobra.Command{
		Use:   "config",
		Short: "Interactively creates a iriscli config file",
		RunE:  runConfigCmd,
	}

	return cfg
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	stdin := BufferStdin()
	gaiaCLIHome, err := handleGaiaCLIHome(home, stdin)
	if err != nil {
		return err
	}
	node, err := handleNode(stdin)
	if err != nil {
		return err
	}
	trustNode, err := handleTrustNode(stdin)
	if err != nil {
		return err
	}

	encoding := "btc"
	output := "text"

	chainID, err := handleChainID(stdin)
	if err != nil {
		return err
	}

	fee, err := handleFee(stdin)
	if err != nil {
		return err
	}

	cfg := &cliConfig{
		Home:      gaiaCLIHome,
		ChainID:   chainID,
		TrustNode: trustNode,
		Encoding:  encoding,
		Output:    output,
		Node:      node,
		Trace:     false,
		Fee:       fee,
	}

	return createGaiaCLIConfig(cfg)
}

func handleGaiaCLIHome(dir string, stdin *bufio.Reader) (string, error) {
	dirName := ".iriscli"
	home, err := GetString(fmt.Sprintf("Where is your iriscli home directory? (Default: ~/%s)", dirName), stdin)
	if err != nil {
		return "", err
	}

	if home == "" {
		home = path.Join(dir, dirName)
	}

	return home, nil
}

func handleNode(stdin *bufio.Reader) (string, error) {
	defaultNode := "tcp://localhost:26657"
	node, err := GetString(fmt.Sprintf("Where is your validator node running? (Default: %s)", defaultNode), stdin)
	if err != nil {
		return "", err
	}

	if node == "" {
		node = defaultNode
	}

	return node, nil
}

func handleChainID(stdin *bufio.Reader) (string, error) {
	chainID, err := GetString("What is your chainID?", stdin)
	if err != nil {
		return "", err
	}

	if chainID == "" {
		return "", fmt.Errorf("you must specify your chainID")
	}

	return chainID, nil
}

func handleFee(stdin *bufio.Reader) (string, error) {
	fee, err := GetString("Please specify default fee", stdin)
	if err != nil {
		return "", err
	}

	if fee == "" {
		fmt.Println("No fee has been specified")
	}

	return fee, nil
}

func handleTrustNode(stdin *bufio.Reader) (bool, error) {
	return GetConfirmation("Do you trust this node?", stdin)
}

func createGaiaCLIConfig(cfg *cliConfig) error {
	cfgPath := path.Join(cfg.Home, "config")
	err := os.MkdirAll(cfgPath, os.ModePerm)
	if err != nil {
		return err
	}

	data, err := toml.Marshal(*cfg)
	if err != nil {
		return err
	}

	cfgFile := path.Join(cfgPath, "config.toml")
	if info, err := os.Stat(cfgFile); err == nil && !info.IsDir() {
		err = os.Rename(cfgFile, path.Join(cfgPath, "config.toml-old"))
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(cfgFile, data, os.ModePerm)
}
