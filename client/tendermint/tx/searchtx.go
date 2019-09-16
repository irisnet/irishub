package tx

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

const (
	flagTags = "tags"
	flagAny  = "any"
	flagPage = "page"
	flagSize = "size"
)

// default client command to search through tagged transactions
func SearchTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs",
		Short: "Search for all transactions that match the given tags.",
		Long: strings.TrimSpace(`
Search for transactions that match exactly the given tags. For example:

$ iriscli tendermint txs --tags '<tag1>:<value1>&<tag2>:<value2>'
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagsStr := viper.GetString(flagTags)
			page := viper.GetInt(flagPage)
			size := viper.GetInt(flagSize)

			tagsStr = strings.Trim(tagsStr, "'")
			var tags []string
			if strings.Contains(tagsStr, "&") {
				tags = strings.Split(tagsStr, "&")
			} else {
				tags = append(tags, tagsStr)
			}

			if page < 0 || size < 0 {
				return fmt.Errorf("page or size should not be negative")
			}

			var tmTags []string
			for _, tag := range tags {
				if !strings.Contains(tag, ":") {
					return fmt.Errorf("%s should be of the format <key>:<value>", tagsStr)
				} else if strings.Count(tag, ":") > 1 {
					return fmt.Errorf("%s should only contain one <key>:<value> pair", tagsStr)
				}
				keyValue := strings.Split(tag, ":")
				if keyValue[0] == types.TxHeightKey {
					tag = fmt.Sprintf("%s=%s", keyValue[0], keyValue[1])
				} else {
					tag = fmt.Sprintf("%s='%s'", keyValue[0], keyValue[1])
				}
				tmTags = append(tmTags, tag)
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txs, err := SearchTxs(cliCtx, cdc, tmTags, page, size)
			if err != nil {
				return err
			}

			var output []byte
			if cliCtx.Indent {
				output, err = cdc.MarshalJSONIndent(txs, "", "  ")
			} else {
				output, err = cdc.MarshalJSON(txs)
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().Bool(client.FlagIndentResponse, true, "Add indent to JSON response")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of Tendermint node")
	cmd.Flags().String(flagTags, "", "tag:value list of tags that must match")
	cmd.Flags().Int(flagPage, 0, "Pagination page")
	cmd.Flags().Int(flagSize, 100, "Pagination size")
	return cmd
}

func SearchTxs(cliCtx context.CLIContext, cdc *codec.Codec, tags []string, page, size int) (*SearchTxsResult, error) {
	if len(tags) == 0 {
		return nil, errors.New("must declare at least one tag to search")
	}

	// XXX: implement ANY
	query := strings.Join(tags, " AND ")

	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	prove := !cliCtx.TrustNode

	res, err := node.TxSearch(query, prove, page, size)
	if err != nil {
		return nil, err
	}

	if prove {
		for _, tx := range res.Txs {
			err := ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, err
			}
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, res.Txs)
	if err != nil {
		return nil, err
	}

	txs, err := FormatTxResults(cdc, res.Txs, resBlocks)
	if err != nil {
		return nil, err
	}

	result := NewSearchTxsResult(res.TotalCount, len(txs), page, size, txs)

	return &result, nil
}

// parse the indexed txs into an array of Info
func FormatTxResults(cdc *codec.Codec, res []*ctypes.ResultTx, resBlocks map[int64]*ctypes.ResultBlock) ([]Info, error) {
	var err error
	out := make([]Info, len(res))
	for i := range res {
		out[i], err = formatTxResult(cdc, res[i], resBlocks[res[i].Height])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// Search Tx REST Handler
func SearchTxRequestHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tags []string
		var txs *SearchTxsResult
		err := r.ParseForm()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("could not parse query parameters", err.Error()))
			return
		}

		if len(r.Form) == 0 {
			utils.PostProcessResponse(w, cdc, txs, cliCtx.Indent)
			return
		}

		for key, values := range r.Form {
			if key == "search_request_page" || key == "search_request_size" {
				continue
			}
			value, err := url.QueryUnescape(values[0])
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("could not decode query value", err.Error()))
				return
			}
			if key == types.TxHeightKey {
				tags = append(tags, fmt.Sprintf("%s=%s", key, value))
			} else {
				tags = append(tags, fmt.Sprintf("%s='%s'", key, value))
			}
		}
		pageString := r.FormValue("search_request_page")
		sizeString := r.FormValue("search_request_size")
		page := int64(0)
		size := int64(100)
		if pageString != "" {
			var ok bool
			page, ok = utils.ParseInt64OrReturnBadRequest(w, pageString)
			if !ok {
				return
			}
		}
		if sizeString != "" {
			var ok bool
			size, ok = utils.ParseInt64OrReturnBadRequest(w, sizeString)
			if !ok {
				return
			}
		}
		if page < 0 || size < 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("page or size should not be negative"))
			return
		}

		txs, err = SearchTxs(cliCtx, cdc, tags, int(page), int(size))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		utils.PostProcessResponse(w, cdc, txs, cliCtx.Indent)
	}
}
