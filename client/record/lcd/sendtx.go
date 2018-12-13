package lcd

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/record"
	sdk "github.com/irisnet/irishub/types"
)

type postRecordReq struct {
	BaseTx      utils.BaseTx `json:"base_tx"`   // basic tx info
	Submitter   string       `json:"submitter"` //  Address of the submitter
	Description string       `json:"description"`
	Data        string       `json:"data"` // for onchain
}

func postRecordHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Init context and read request parameters
		cliCtx = utils.InitReqCliCtx(cliCtx, r)

		var req postRecordReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		submitter, err := sdk.AccAddressFromBech32(req.Submitter)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		onchainData := req.Data

		var recordHash string
		var dataSize int64
		// --onchain-data has a high priority over --file-path
		if len(onchainData) != 0 {
			dataSize = int64(binary.Size([]byte(onchainData)))
			if dataSize >= record.UploadLimitOfOnchain {
				utils.WriteErrorResponse(w, http.StatusBadRequest,
					fmt.Sprintf("Upload data is too large, max supported data size is %d", record.UploadLimitOfOnchain))
				return
			}

			sum := sha256.Sum256([]byte(onchainData))
			recordHash = hex.EncodeToString(sum[:])
		} else {
			utils.WriteErrorResponse(w, http.StatusBadRequest,
				"--onchain-data is empty and pleae double check this option")
			return
		}

		recordID := record.KeyRecord(recordHash)
		res, err := cliCtx.QueryStore([]byte(recordID), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) != 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Record ID '%s' is already existed", recordID))
			return
		}

		submitTime := time.Now().Unix()

		// create the message
		msg := record.NewMsgSubmitRecord(req.Description, submitTime, submitter, recordHash, dataSize, onchainData)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
