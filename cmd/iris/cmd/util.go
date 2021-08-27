package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/olebedev/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

const (
	formatJSON     = "json"
	cmdScopeGlobal = "global"
	filedTypeMap   = "map"
	filedTypeArray = "array"
)

var (
	converter = NewConverter().
			registerGlobalFlag("fees").
			registerGlobalFlag("amount").
			registerGlobalFlag("deposit").
			registerGlobalFlag("service-fee-cap").
			registerCmdWithArgs("bank", "send", 2).
			registerCmdWithArgs("staking", "delegate", 1).
			registerCmdWithArgs("staking", "redelegate", 2).
			registerCmdWithArgs("staking", "unbond", 1).
			registerCmdWithArgs("distribution", "fund-community-pool", 0).
			registerCmdWithArgs("gov", "deposit", 1).
			registerCmdWithArgs("ibc-transfer", "transfer", 3).
			registerCmdForResponse("bank", "balances", "balances", filedTypeArray).
			registerCmdForResponse("bank", "total", "supply", filedTypeArray).
			registerCmdForResponse("gov", "params", "deposit_params.min_deposit", filedTypeArray).
			registerCmdForResponse("distribution", "validator-outstanding-rewards", "rewards", filedTypeArray).
			registerCmdForResponses("distribution", "rewards",
			field{name: "total", typ: filedTypeArray},
			field{name: "rewards.*.reward", typ: filedTypeArray},
			field{name: "rewards", typ: filedTypeArray}).
		registerCmdForResponse("token", "total-burn", "burned_coins", filedTypeArray)

	rescueStdout = os.Stdout
)

type (
	field struct {
		name  string
		index int
		typ   string
	}

	command struct {
		parentCmd string
		fields    map[string]field
	}
)

func (c command) append(name, typ string, index int) command {
	c.fields[name] = field{
		name:  name,
		index: index,
		typ:   typ,
	}
	return c
}

type coinConverter struct {
	cmds   map[string]command
	tokens map[string]tokentypes.TokenI
	r, w   *os.File
}

// NewConverter return a instance of coinConverter
func NewConverter() *coinConverter {
	return &coinConverter{
		cmds:   make(map[string]command),
		tokens: make(map[string]tokentypes.TokenI),
	}
}

func (it *coinConverter) key(parentCmd, cmd string) string {
	return fmt.Sprintf("%s/%s", parentCmd, cmd)
}

func (it *coinConverter) registerCmdWithArgs(parentCmd, cmd string, argsIdx int) *coinConverter {
	commands, ok := it.cmds[cmd]
	if !ok {
		commands = command{
			parentCmd: parentCmd,
			fields:    map[string]field{},
		}
	}
	commands = commands.append("ARGS", "", argsIdx)
	it.cmds[cmd] = commands
	return it
}

func (it *coinConverter) registerGlobalFlag(flagNm string) *coinConverter {
	commands, ok := it.cmds[cmdScopeGlobal]
	if !ok {
		commands = command{
			parentCmd: cmdScopeGlobal,
			fields:    map[string]field{},
		}
	}
	commands = commands.append(flagNm, "", -1)
	it.cmds[cmdScopeGlobal] = commands
	return it
}

func (it *coinConverter) registerCmdWithFlag(parentCmd, cmd, flagNm string) *coinConverter {
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

func (it *coinConverter) registerCmdForResponse(parentCmd, cmd, jsonPath, typ string) *coinConverter {
	commands, ok := it.cmds[cmd]
	if !ok {
		commands = command{
			parentCmd: parentCmd,
			fields:    map[string]field{},
		}
	}

	commands = commands.append(jsonPath, typ, -1)

	it.cmds[cmd] = commands
	return it
}

func (it *coinConverter) registerCmdForResponses(parentCmd, cmd string, fields ...field) *coinConverter {
	commands, ok := it.cmds[cmd]
	if !ok {
		commands = command{
			parentCmd: parentCmd,
			fields:    map[string]field{},
		}
	}

	for _, field := range fields {
		commands = commands.append(field.name, field.typ, -1)
	}

	it.cmds[cmd] = commands
	return it
}

func (it coinConverter) hasFromFlag(cmdNm, flagNm string) bool {
	if cmd, ok := it.cmds[cmdNm]; ok {
		if _, ok = cmd.fields[flagNm]; ok {
			return ok
		}
	}
	cmd := it.cmds[cmdScopeGlobal]
	_, ok := cmd.fields[flagNm]
	return ok
}

func (it coinConverter) getFromArgs(cmdNm string) (field, bool) {
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

func (it coinConverter) getFields(cmdNm string) map[string]field {
	cmd, ok := it.cmds[cmdNm]
	if !ok {
		return map[string]field{}
	}
	return cmd.fields
}

func (it *coinConverter) handlePreRun(cmd *cobra.Command, args []string) {
	if b, _ := cmd.Flags().GetBool(flags.FlagGenerateOnly); b {
		return
	}
	cmdNm := cmd.Name()
	//handle flag
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed {
			viper.SetDefault(flag.Name, flag.Value)
		}
		it.parseFlags(cmd, flag, cmdNm)
	})

	//handle field
	it.parseArgs(cmd, args[:])

	if !it.isOutputYAML(cmd) {
		return
	}
	it.r, it.w, _ = os.Pipe()
	os.Stdout = it.w
}

func (it *coinConverter) handlePostRun(cmd *cobra.Command) {
	if !it.isOutputYAML(cmd) {
		return
	}
	if it.w != nil {
		_ = it.w.Close()
	}
	out, _ := ioutil.ReadAll(it.r)
	os.Stdout = rescueStdout
	fmt.Println(it.parseYAML(cmd, out))
}

func (it coinConverter) parseFlags(cmd *cobra.Command, flag *pflag.Flag, cmdNm string) {
	if it.hasFromFlag(cmdNm, flag.Name) {
		srcCoinStr := flag.Value.String()
		if res, err := it.convertCoins(cmd, srcCoinStr); err == nil {
			_ = flag.Value.Set(res)
		}
	}
}

func (it coinConverter) parseArgs(cmd *cobra.Command, args []string) {
	command, ok := it.cmds[cmd.Name()]
	if !ok {
		return
	}

	if cmd.Parent().Name() != command.parentCmd {
		return
	}

	if field, ok := it.getFromArgs(cmd.Name()); ok && len(args) > 0 {
		res, err := it.convertCoins(cmd, args[field.index])
		if err != nil {
			return
		}
		args[field.index] = res
	}
}

func (it coinConverter) parseYAML(cmd *cobra.Command, in []byte) string {
	cfg, err := config.ParseYamlBytes(in)
	if err != nil {
		return string(in)
	}

	fields := it.getFields(cmd.Name())
	for path, field := range fields {
		for _, p := range it.resolvePath(cfg, path) {
			switch field.typ {
			case filedTypeArray:
				it.handleList(cmd, cfg, p)
			case filedTypeMap:
				it.handleMap(cmd, cfg, p)
			}
		}
	}
	s, err := config.RenderYaml(cfg.Root)
	if err != nil {
		return string(in)
	}
	return s
}

func (it coinConverter) resolvePath(cfg *config.Config, path string) (paths []string) {
	subPaths := strings.SplitN(path, "*", 2)
	if len(subPaths) == 1 {
		return []string{path}
	}

	if len(subPaths) != 2 {
		return paths
	}

	list, err := cfg.List(subPaths[0][:len(subPaths[0])-1])
	if err != nil {
		return paths
	}

	for i := 0; i < len(list); i++ {
		p := fmt.Sprintf("%s%d%s", subPaths[0], i, subPaths[1])
		paths = append(paths, it.resolvePath(cfg, p)...)
	}
	return paths
}

func (it *coinConverter) queryToken(cmd *cobra.Command, denom string) (ft tokentypes.TokenI, err error) {
	if ft, ok := it.tokens[denom]; ok {
		return ft, nil
	}

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return nil, err
	}

	if err := tokentypes.ValidateSymbol(denom); err != nil {
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

	it.tokens[denom] = evi
	return evi, nil
}

func (it *coinConverter) isOutputYAML(cmd *cobra.Command) bool {
	output1, err := cmd.Flags().GetString(cli.OutputFlag)
	output2 := viper.GetString(cli.OutputFlag)
	if output2 == formatJSON || (err == nil && output1 == formatJSON) {
		return false
	}
	cmdPath := cmd.CommandPath()
	if !strings.Contains(cmdPath, queryCommand().CommandPath()) {
		return false
	}
	return true
}

func (it *coinConverter) handleList(cmd *cobra.Command, cfg *config.Config, path string) {
	list, err := cfg.List(path)
	if err != nil {
		return
	}
	for i := range list {
		it.handleMap(cmd, cfg, fmt.Sprintf("%s.%d", path, i))
	}
}

func (it *coinConverter) handleMap(cmd *cobra.Command, cfg *config.Config, path string) {
	cMap, err := cfg.Map(path)
	if err != nil {
		return
	}

	bz, err := json.Marshal(cMap)
	if err != nil {
		return
	}

	var srcCoin sdk.DecCoin
	if err := json.Unmarshal(bz, &srcCoin); err != nil || !srcCoin.IsValid() {
		return
	}

	truncCoin, _ := srcCoin.TruncateDecimal()
	dstCoin, err := it.convertToMainCoin(cmd, truncCoin)
	if err != nil {
		return
	}
	_ = cfg.Set(path, dstCoin)
}

func (it *coinConverter) convertCoins(cmd *cobra.Command, coinsStr string) (dstCoinsStr string, err error) {
	cs, err := it.parseCoins(coinsStr)
	if err != nil {
		return coinsStr, err
	}
	dstCoins := sdk.Coins{}
	for _, coin := range cs {
		if c, err := it.convertToMinCoin(cmd, coin); err == nil {
			dstCoins = append(dstCoins, c)
			continue
		}
		c, _ := coin.TruncateDecimal()
		dstCoins = append(dstCoins, c)
	}
	return dstCoins.String(), nil
}

func (it *coinConverter) convertToMinCoin(cmd *cobra.Command, srcCoin sdk.DecCoin) (coin sdk.Coin, err error) {
	ft, err := it.queryToken(cmd, srcCoin.Denom)
	if err != nil {
		return coin, err
	}
	return ft.ToMinCoin(srcCoin)
}

func (it *coinConverter) convertToMainCoin(cmd *cobra.Command, srcCoin sdk.Coin) (coin sdk.DecCoin, err error) {
	ft, err := it.queryToken(cmd, srcCoin.Denom)
	if err != nil {
		return coin, err
	}
	return ft.ToMainCoin(srcCoin)
}

func (it *coinConverter) parseCoins(srcCoinsStr string) (sdk.DecCoins, error) {
	if cs, err := sdk.ParseDecCoins(srcCoinsStr); err == nil {
		return cs, nil
	}
	if cs, err := sdk.ParseCoinsNormalized(srcCoinsStr); err == nil {
		return sdk.NewDecCoinsFromCoins(cs...), nil
	}
	return sdk.DecCoins{}, fmt.Errorf("parsed decimal coins are invalid: %s", srcCoinsStr)
}
