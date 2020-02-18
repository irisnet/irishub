package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/oracle"
	"github.com/irisnet/irishub/client/context"
	cli "github.com/irisnet/irishub/client/oracle"
	oracleClient "github.com/irisnet/irishub/client/oracle"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// query a feed definition
	r.HandleFunc(
		fmt.Sprintf("/oracle/feeds/{%s}", FeedName),
		queryFeedHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a feed list by condition
	r.HandleFunc(
		fmt.Sprintf("/oracle/feeds"),
		queryFeedsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a feed value by feed name
	r.HandleFunc(
		fmt.Sprintf("/oracle/feeds/{%s}/values", FeedName),
		queryFeedValuesHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

func queryFeedHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		feedName := vars[FeedName]

		params := oracle.QueryFeedParams{
			FeedName: feedName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.OracleRoute, oracle.QueryFeed)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryFeedsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := oracle.QueryFeedsParams{
			State: cli.GetState(oracleClient.GetUrlParam(r.URL, FeedState)),
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.OracleRoute, oracle.QueryFeeds)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryFeedValuesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := oracle.QueryFeedsParams{
			State: cli.GetState(oracleClient.GetUrlParam(r.URL, FeedState)),
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.OracleRoute, oracle.QueryFeedValue)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}
