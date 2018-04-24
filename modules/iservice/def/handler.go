package def

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	cmn "github.com/tendermint/tmlibs/common"
	"log"
	"strings"
)

type SvcDefService struct {
	keeper SvcDefKeeper
}

func NewSvcDefService(keeper SvcDefKeeper) SvcDefService {
	return SvcDefService{
		keeper,
	}
}

func (service SvcDefService) CheckTx(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	log.Printf("enter into %s", "SvcDefService.CheckTx")

	m := msg.(SvcDefMsg)

	if service.keeper.Has(ctx, GetSvcDefKey(m.Name)) {
		sdkErr := HasExisted("")
		return sdkErr.Result()
	}
	return sdk.Result{}
}

func (service SvcDefService) DeliverTx(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	m := msg.(SvcDefMsg)
	tags := make([]cmn.KVPair, 0, 5)
	tags = append(tags, IndexHeight(ctx.BlockHeight()))
	tags = append(tags, IndexServiceName(m.Name))
	tags = append(tags, IndexChainId(m.ChainID))
	tags = append(tags, IndexMessagingType(m.Messaging))

	addr := fmt.Sprintf("%s", m.Creator)
	tags = append(tags, IndexSender(addr))

	if m.Tags != "" {
		tx_tags := strings.Split(m.Tags, ",")
		for _, tag := range tx_tags {
			kv := strings.Split(tag, "=")
			if len(kv) == 2 {
				tags = append(tags, IndexKVTag(kv[0], kv[1]))
			} else if len(kv) == 1 && kv[0] != "" {
				tags = append(tags, IndexKeyTag(kv[0]))
			}
		}
	}

	cdc := wire.NewCodec()
	bz, _ := cdc.MarshalBinary(msg)

	service.keeper.Set(ctx, GetSvcDefKey(m.Name), bz)
	return sdk.Result{
		Tags: tags,
	}
}
