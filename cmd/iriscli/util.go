package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/olebedev/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

const (
	formatJson      = "json"
	tokenQueryRoute = "asset/token"

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

func handleRequestPreRun(cdc *codec.Codec, cmd *cobra.Command, args []string) {
	if b, _ := cmd.Flags().GetBool(flags.FlagGenerateOnly); b {
		return
	}
	cmdNm := cmd.Name()
	//handle flag
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed {
			viper.SetDefault(flag.Name, flag.Value)
		}
		parseFlags(cdc, flag, cmdNm)
	})

	//handle field
	parseArgs(cdc, cmdNm, args[:])
}

func handleResponsePreRun(cdc *codec.Codec, cmd *cobra.Command) {
	if !isOutputYAML(cdc, cmd) {
		return
	}
	r, w, _ = os.Pipe()
	os.Stdout = w
}

func handleResponsePostRun(cdc *codec.Codec, cmd *cobra.Command) {
	if !isOutputYAML(cdc, cmd) {
		return
	}
	if w != nil {
		_ = w.Close()
	}
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	fmt.Println(parseYAML(cdc, cmd.Name(), out))
}

func isOutputYAML(cdc *codec.Codec, cmd *cobra.Command) bool {
	output1, err := cmd.Flags().GetString(cli.OutputFlag)
	output2 := viper.GetString(cli.OutputFlag)
	if output2 == formatJson || (err == nil && output1 == formatJson) {
		return false
	}
	cmdPath := cmd.CommandPath()
	if !strings.Contains(cmdPath, queryCmd(cdc).CommandPath()) {
		return false
	}
	return true
}

func parseFlags(cdc *codec.Codec, flag *pflag.Flag, cmdNm string) {
	if cmdCfg.hasFromFlag(cmdNm, flag.Name) {
		srcCoinStr := flag.Value.String()
		if res, err := convertCoins(cdc, srcCoinStr); err == nil {
			_ = flag.Value.Set(res)
		}
	}
}

func parseArgs(cdc *codec.Codec, cmd string, args []string) {
	if field, ok := cmdCfg.getFromArgs(cmd); ok && len(args) > 0 {
		if res, err := convertCoins(cdc, args[field.index]); err == nil {
			args[field.index] = res
		}
	}
}

func parseYAML(cdc *codec.Codec, cmd string, in []byte) string {
	srcYaml := string(in)
	dstStr := strings.ReplaceAll(srcYaml, "|", "")
	cfg, err := config.ParseYaml(dstStr)
	if err != nil {
		return srcYaml
	}
	for k, v := range cmdCfg.getFromResponse(cmd) {
		switch v.typ {
		case filedTypeArray:
			handleList(cdc, cfg, k)
		case filedTypeMap:
			handleMap(cdc, cfg, k)
		}
	}
	s, err := config.RenderYaml(cfg.Root)
	if err != nil {
		return srcYaml
	}
	return s
}

func handleList(cdc *codec.Codec, cfg *config.Config, path string) {
	list, err := cfg.List(path)
	if err != nil {
		return
	}
	for i, _ := range list {
		handleMap(cdc, cfg, fmt.Sprintf("%s.%d", path, i))
	}
}

func handleMap(cdc *codec.Codec, cfg *config.Config, path string) {
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
	dstCoin, err := convertToMainCoin(cdc, srcCoin)
	if err != nil {
		return
	}
	_ = cfg.Set(path, dstCoin)
}

func convertCoins(cdc *codec.Codec, coinsStr string) (dstCoinsStr string, err error) {
	cs, err := parseCoins(coinsStr)
	if err != nil {
		return coinsStr, err
	}
	dstCoins := sdk.Coins{}
	for _, coin := range cs {
		if c, err := convertToMinCoin(cdc, coin); err == nil {
			dstCoins = append(dstCoins, c)
			continue
		}
		c, _ := coin.TruncateDecimal()
		dstCoins = append(dstCoins, c)
	}
	return dstCoins.String(), nil
}

func convertToMinCoin(cdc *codec.Codec, srcCoin sdk.DecCoin) (coin sdk.Coin, err error) {
	params := token.QueryTokenParams{
		Symbol: srcCoin.Denom,
	}

	ft, err := queryToken(cdc, params)
	if err != nil {
		return coin, err
	}
	return ft.ToMinCoin(srcCoin)
}

func convertToMainCoin(cdc *codec.Codec, srcCoin sdk.Coin) (coin sdk.DecCoin, err error) {
	params := token.QueryTokenParams{
		MinUnit: srcCoin.Denom,
	}

	ft, err := queryToken(cdc, params)
	if err != nil {
		return coin, err
	}
	return ft.ToMainCoin(srcCoin)
}

func queryToken(cdc *codec.Codec, params token.QueryTokenParams) (ft token.FungibleToken, err error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return ft, err
	}

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", tokenQueryRoute, token.QueryToken), bz,
	)

	if err != nil {
		return ft, err
	}

	var fts token.Tokens
	cliCtx.Codec.MustUnmarshalJSON(res, &fts)
	if len(fts) == 0 {
		return
	}
	return fts[0], nil
}

func parseCoins(srcCoinsStr string) (sdk.DecCoins, error) {
	if cs, err := sdk.ParseDecCoins(srcCoinsStr); err == nil {
		return cs, nil
	}
	if cs, err := sdk.ParseCoins(srcCoinsStr); err == nil {
		return sdk.NewDecCoins(cs), nil
	}
	return sdk.DecCoins{}, fmt.Errorf("parsed decimal coins are invalid: %s", srcCoinsStr)
}
