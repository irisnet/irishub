package rpc

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/tendermint/tendermint/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmprotoversion "github.com/tendermint/tendermint/proto/tendermint/version"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type Protocol uint64
type Address = crypto.Address
type SignedMsgType byte

// Single block (with meta)
type ResultBlock struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}

// Block defines the atomic unit of a Tendermint blockchain.
type Block struct {
	mtx          sync.Mutex
	Header       `json:"header"`
	tmtypes.Data `json:"data"`
	Evidence     EvidenceData `json:"evidence"`
	LastCommit   *Commit      `json:"last_commit"`
}

// BlockMeta contains meta information about a block - namely, it's ID and Header.
type BlockMeta struct {
	BlockID BlockID `json:"block_id"` // the block hash and partsethash
	Header  Header  `json:"header"`   // The block's Header
}

type BlockID struct {
	Hash          tmbytes.HexBytes `json:"hash"`
	PartSetHeader PartSetHeader    `json:"parts"`
}

type PartSetHeader struct {
	Total int              `json:"total"`
	Hash  tmbytes.HexBytes `json:"hash"`
}

// MaxDataBytesUnknownEvide
// Header defines the structure of a Tendermint block header.
// NOTE: changes to the Header should be duplicated in:
// - header.Hash()
// - abci.Header
// - /docs/spec/blockchain/blockchain.md
type Header struct {
	// basic block info
	Version  tmprotoversion.Consensus `json:"version"`
	ChainID  string                   `json:"chain_id"`
	Height   int64                    `json:"height"`
	Time     time.Time                `json:"time"`
	NumTxs   int64                    `json:"num_txs"`
	TotalTxs int64                    `json:"total_txs"`

	// prev block info
	LastBlockID BlockID `json:"last_block_id"`

	// hashes of block data
	LastCommitHash tmbytes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       tmbytes.HexBytes `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     tmbytes.HexBytes `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash tmbytes.HexBytes `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      tmbytes.HexBytes `json:"consensus_hash"`       // consensus params for current block
	AppHash            tmbytes.HexBytes `json:"app_hash"`             // state after txs from the previous block
	LastResultsHash    tmbytes.HexBytes `json:"last_results_hash"`    // root hash of all results from the txs from the previous block

	// consensus info
	EvidenceHash    tmbytes.HexBytes `json:"evidence_hash"`    // evidence included in the block
	ProposerAddress Address          `json:"proposer_address"` // original proposer of the block
}

// Commit contains the evidence that a block was committed by a set of validators.
// NOTE: Commit is empty for height 1, but never nil.
type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	BlockID    BlockID `json:"block_id"`
	Precommits []*Vote `json:"precommits"`

	// Volatile
	firstPrecommit *Vote
	hash           tmbytes.HexBytes
	bitArray       *BitArray
}

type Vote struct {
	Type             SignedMsgType `json:"type"`
	Height           int64         `json:"height"`
	Round            int32         `json:"round"`
	BlockID          BlockID       `json:"block_id"` // zero if vote is nil.
	Timestamp        time.Time     `json:"timestamp"`
	ValidatorAddress Address       `json:"validator_address"`
	ValidatorIndex   int           `json:"validator_index"`
	Signature        []byte        `json:"signature"`
}

// BitArray is a thread-safe implementation of a bit array.
type BitArray struct {
	mtx   sync.Mutex
	Bits  int      `json:"bits"`  // NOTE: persisted via reflect, must be exported
	Elems []uint64 `json:"elems"` // NOTE: persisted via reflect, must be exported
}

// SignedHeader is a header along with the commits that prove it.
// It is the basis of the lite client.
type SignedHeader struct {
	*Header `json:"header"`
	Commit  *Commit `json:"commit"`
}

//-----------------------------------------------------------------------------

// EvidenceData contains any evidence of malicious wrong-doing by validators
type EvidenceData struct {
	Evidence tmtypes.EvidenceList `json:"evidence"`

	// Volatile
	hash tmbytes.HexBytes
}

func getBlock(clientCtx client.Context, height *int64) ([]byte, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(context.Background(), height)
	if err != nil {
		return nil, err
	}

	result := convertResultBlock(res)

	return legacy.Cdc.MarshalJSON(result)
}

func convertResultBlock(tmBlock *ctypes.ResultBlock) *ResultBlock {
	precommits := make([]*Vote, len(tmBlock.Block.LastCommit.Signatures))
	lastCommitBlockID := BlockID{
		Hash: tmBlock.Block.LastCommit.BlockID.Hash,
		PartSetHeader: PartSetHeader{
			Total: int(tmBlock.Block.LastCommit.BlockID.PartSetHeader.Total),
			Hash:  tmBlock.Block.LastCommit.BlockID.PartSetHeader.Hash,
		},
	}
	lastBlockID := BlockID{
		Hash: tmBlock.Block.LastBlockID.Hash,
		PartSetHeader: PartSetHeader{
			Total: int(tmBlock.Block.LastBlockID.PartSetHeader.Total),
			Hash:  tmBlock.Block.LastBlockID.PartSetHeader.Hash,
		},
	}
	blockID := BlockID{
		Hash: tmBlock.BlockID.Hash,
		PartSetHeader: PartSetHeader{
			Total: int(tmBlock.BlockID.PartSetHeader.Total),
			Hash:  tmBlock.BlockID.PartSetHeader.Hash,
		},
	}
	for k, v := range tmBlock.Block.LastCommit.Signatures {
		precommits[k] = &Vote{
			Type:             0x02,
			Height:           tmBlock.Block.LastCommit.Height,
			Round:            tmBlock.Block.LastCommit.Round,
			BlockID:          lastCommitBlockID,
			Timestamp:        v.Timestamp,
			ValidatorAddress: v.ValidatorAddress,
			ValidatorIndex:   k,
			Signature:        v.Signature,
		}
	}
	header := Header{
		Version:            tmBlock.Block.Version,
		ChainID:            tmBlock.Block.ChainID,
		Height:             tmBlock.Block.Height,
		Time:               tmBlock.Block.Time,
		NumTxs:             0,
		TotalTxs:           int64(len(tmBlock.Block.Data.Txs)),
		LastBlockID:        lastBlockID,
		LastCommitHash:     tmBlock.Block.LastCommitHash,
		DataHash:           tmBlock.Block.DataHash,
		ValidatorsHash:     tmBlock.Block.ValidatorsHash,
		NextValidatorsHash: tmBlock.Block.NextValidatorsHash,
		ConsensusHash:      tmBlock.Block.ConsensusHash,
		AppHash:            tmBlock.Block.AppHash,
		LastResultsHash:    tmBlock.Block.LastResultsHash,
		EvidenceHash:       tmBlock.Block.EvidenceHash,
		ProposerAddress:    tmBlock.Block.ProposerAddress,
	}
	return &ResultBlock{
		BlockMeta: &BlockMeta{
			BlockID: blockID,
			Header:  header,
		},
		Block: &Block{
			Header: header,
			Data:   tmBlock.Block.Data,
			Evidence: EvidenceData{
				Evidence: tmBlock.Block.Evidence.Evidence,
			},
			LastCommit: &Commit{
				BlockID:    lastCommitBlockID,
				Precommits: precommits,
			},
		},
	}
}

// get the current blockchain height
func GetChainHeight(clientCtx client.Context) (int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return -1, err
	}

	status, err := node.Status(context.Background())
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// REST handler to get a block
func BlockRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}

		chainHeight, err := GetChainHeight(clientCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlock(clientCtx, &height)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}

// REST handler to get the latest block
func LatestBlockRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlock(clientCtx, nil)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponseBare(w, clientCtx, output)
	}
}
