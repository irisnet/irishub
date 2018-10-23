package rpc

import (
	"fmt"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"net/http"
	"strconv"
)

func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query remote node for status",
		Example: "iriscli status",
		RunE:  printNodeStatus,
	}

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	return cmd
}

func GetNodeStatus(cliCtx context.CLIContext) (*ctypes.ResultStatus, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}

	return node.Status()
}

// CMD

func printNodeStatus(cmd *cobra.Command, args []string) error {
	status, err := GetNodeStatus(context.NewCLIContext())
	if err != nil {
		return err
	}

	output, err := cdc.MarshalJSON(status)
	// output, err := cdc.MarshalJSONIndent(res, "  ", "")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

// REST handler for node info
func NodeInfoRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := GetNodeStatus(cliCtx)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		nodeInfo := status.NodeInfo
		output, err := cdc.MarshalJSONIndent(nodeInfo,"", "  ")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// REST handler for node syncing
func NodeSyncingRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := GetNodeStatus(cliCtx)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		syncing := status.SyncInfo.CatchingUp
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(strconv.FormatBool(syncing)))
	}
}
