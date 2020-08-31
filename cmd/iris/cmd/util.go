package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tokentypes "github.com/irismod/token/types"
	"github.com/olebedev/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	formatJson = "json"

	cmdScopeGlobal = "global"

	filedTypeMap   = "map"
	filedTypeArray = "array"
)

var (
	cmdCfg = NewCmdConfig().appendFromFlag(cmdScopeGlobal, cmdScopeGlobal, "fees").
		appendFromFlag(cmdScopeGlobal, cmdScopeGlobal, "amount").
		appendFromArgs("bank", "send", 2).
		appendFromArgs("staking", "delegate", 1).
		appendFromArgs("staking", "redelegate", 2).
		appendFromArgs("staking", "unbond", 1).
		appendFromArgs("distribution", "fund-community-pool", 0).
		appendFromArgs("token", "burn", 0).
		appendFromArgs("gov", "deposit", 1).
		appendFromResponse("query", "account", "coins", filedTypeArray).
		appendFromResponse("token", "params", "issue_token_base_fee", filedTypeMap).
		appendFromResponse("gov", "params", "deposit_parmas.min_deposit", filedTypeArray)

	rescueStdout = os.Stdout
	r, w         *os.File
)

type field struct {
	name  string
	index int
	typ   string
}

type command struct {
	parentCmd string
	fields    map[string]field
}

func (c command) append(name, typ string, index int) command {
	c.fields[name] = field{
		name:  name,
		index: index,
		typ:   typ,
	}
	return c
}

type cmdConfig struct {
	cmds map[string]command
}

func NewCmdConfig() *cmdConfig {
	return &cmdConfig{
		cmds: map[string]command{},
	}
}

func (it *cmdConfig) key(parentCmd, cmd string) string {
	return fmt.Sprintf("%s/%s", parentCmd, cmd)
}

func (it *cmdConfig) appendFromArgs(parentCmd, cmd string, index int) *cmdConfig {
	commands, ok := it.cmds[cmd]
	if !ok {
		commands = command{
			parentCmd: parentCmd,
			fields:    map[string]field{},
		}
	}
	commands = commands.append("ARGS", "", index)
	it.cmds[cmd] = commands
	return it
}

func (it *cmdConfig) appendFromFlag(parentCmd, cmd, flagNm string) *cmdConfig {
	commands, ok := it.cmds[cmd]
	if !ok {
		commands = command{
			parentCmd: parentCmd,
			fields:    map[string]field{},
		}
	}
	commands = commands.append(flagNm, "", -1)
	it.cmds[cmd] = commands
	return it
}

func (it *cmdConfig) appendFromResponse(parentCmd, cmd, flagNm, typ string) *cmdConfig {
	commands, ok := it.cmds[cmd]
	if !ok {
		commands = command{
			parentCmd: parentCmd,
			fields:    map[string]field{},
		}
	}
	commands = commands.append(flagNm, typ, -1)
	it.cmds[cmd] = commands
	return it
}

func (it cmdConfig) hasFromFlag(cmdNm, flagNm string) bool {
	if cmd, ok := it.cmds[cmdNm]; ok {
		if _, ok = cmd.fields[flagNm]; ok {
			return ok
		}
	}
	cmd := it.cmds[cmdScopeGlobal]
	_, ok := cmd.fields[flagNm]
	return ok
}

func (it cmdConfig) getFromArgs(cmdNm string) (field, bool) {
	cmd, ok := it.cmds[cmdNm]
	if !ok {
		return field{}, false
	}
	f := cmd.fields["ARGS"]
	if f.index < 0 {
		return f, false
	}
	return cmd.fields["ARGS"], true
}

func (it cmdConfig) getFromResponse(cmdNm string) map[string]field {
	cmd, ok := it.cmds[cmdNm]
	if !ok {
		return map[string]field{}
	}
	return cmd.fields
}

func handleRequestPreRun(cmd *cobra.Command, args []string) {
	if b, _ := cmd.Flags().GetBool(flags.FlagGenerateOnly); b {
		return
	}
	cmdNm := cmd.Name()
	//handle flag
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed {
			viper.SetDefault(flag.Name, flag.Value)
		}
		parseFlags(cmd, flag, cmdNm)
	})

	//handle field
	parseArgs(cmd, args[:])
}

func handleResponsePreRun(cmd *cobra.Command) {
	if !isOutputYAML(cmd) {
		return
	}
	r, w, _ = os.Pipe()
	os.Stdout = w
}

func handleResponsePostRun(cdc codec.JSONMarshaler, cmd *cobra.Command) {
	if !isOutputYAML(cmd) {
		return
	}
	if w != nil {
		_ = w.Close()
	}
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	fmt.Println(parseYAML(cmd, out))
}

func isOutputYAML(cmd *cobra.Command) bool {
	output1, err := cmd.Flags().GetString(cli.OutputFlag)
	output2 := viper.GetString(cli.OutputFlag)
	if output2 == formatJson || (err == nil && output1 == formatJson) {
		return false
	}
	cmdPath := cmd.CommandPath()
	if !strings.Contains(cmdPath, queryCommand().CommandPath()) {
		return false
	}
	return true
}

func parseFlags(cmd *cobra.Command, flag *pflag.Flag, cmdNm string) {
	if cmdCfg.hasFromFlag(cmdNm, flag.Name) {
		srcCoinStr := flag.Value.String()
		if res, err := convertCoins(cmd, srcCoinStr); err == nil {
			_ = flag.Value.Set(res)
		}
	}
}

func parseArgs(cmd *cobra.Command, args []string) {
	if field, ok := cmdCfg.getFromArgs(cmd.Name()); ok && len(args) > 0 {
		if res, err := convertCoins(cmd, args[field.index]); err == nil {
			args[field.index] = res
		}
	}
}

func parseYAML(cmd *cobra.Command, in []byte) string {
	srcYaml := string(in)
	dstStr := strings.ReplaceAll(srcYaml, "|", "")
	cfg, err := config.ParseYaml(dstStr)
	if err != nil {
		return srcYaml
	}
	for k, v := range cmdCfg.getFromResponse(cmd.Name()) {
		switch v.typ {
		case filedTypeArray:
			handleList(cmd, cfg, k)
		case filedTypeMap:
			handleMap(cmd, cfg, k)
		}
	}
	s, err := config.RenderYaml(cfg.Root)
	if err != nil {
		return srcYaml
	}
	return s
}

func handleList(cmd *cobra.Command, cfg *config.Config, path string) {
	list, err := cfg.List(path)
	if err != nil {
		return
	}
	for i := range list {
		handleMap(cmd, cfg, fmt.Sprintf("%s.%d", path, i))
	}
}

func handleMap(cmd *cobra.Command, cfg *config.Config, path string) {
	cMap, err := cfg.Map(path)
	if err != nil {
		return
	}
	bz, err := json.Marshal(cMap)
	if err != nil {
		return
	}
	var srcCoin sdk.Coin
	if err := json.Unmarshal(bz, &srcCoin); err != nil {
		return
	}
	dstCoin, err := convertToMainCoin(cmd, srcCoin)
	if err != nil {
		return
	}
	_ = cfg.Set(path, dstCoin)
}

func convertCoins(cmd *cobra.Command, coinsStr string) (dstCoinsStr string, err error) {
	cs, err := parseCoins(coinsStr)
	if err != nil {
		return coinsStr, err
	}
	dstCoins := sdk.Coins{}
	for _, coin := range cs {
		if c, err := convertToMinCoin(cmd, coin); err == nil {
			dstCoins = append(dstCoins, c)
			continue
		}
		c, _ := coin.TruncateDecimal()
		dstCoins = append(dstCoins, c)
	}
	return dstCoins.String(), nil
}

func convertToMinCoin(cmd *cobra.Command, srcCoin sdk.DecCoin) (coin sdk.Coin, err error) {
	ft, err := queryToken(cmd, srcCoin.Denom)
	if err != nil {
		return coin, err
	}
	return ft.ToMinCoin(srcCoin)
}

func convertToMainCoin(cmd *cobra.Command, srcCoin sdk.Coin) (coin sdk.DecCoin, err error) {
	ft, err := queryToken(cmd, srcCoin.Denom)
	if err != nil {
		return coin, err
	}
	return ft.ToMainCoin(srcCoin)
}

func queryToken(cmd *cobra.Command, denom string) (ft tokentypes.TokenI, err error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	clientCtx, err = client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
	if err != nil {
		return nil, err
	}

	if err := tokentypes.CheckSymbol(denom); err != nil {
		return nil, err
	}

	queryClient := tokentypes.NewQueryClient(clientCtx)

	res, err := queryClient.Token(context.Background(), &tokentypes.QueryTokenRequest{
		Denom: denom,
	})
	if err != nil {
		return nil, err
	}

	var evi tokentypes.TokenI
	err = clientCtx.InterfaceRegistry.UnpackAny(res.Token, &evi)
	if err != nil {
		return nil, err
	}

	return evi, nil
}

func parseCoins(srcCoinsStr string) (sdk.DecCoins, error) {
	if cs, err := sdk.ParseDecCoins(srcCoinsStr); err == nil {
		return cs, nil
	}
	if cs, err := sdk.ParseCoins(srcCoinsStr); err == nil {
		return sdk.NewDecCoinsFromCoins(cs...), nil
	}
	return sdk.DecCoins{}, fmt.Errorf("parsed decimal coins are invalid: %s", srcCoinsStr)
}
