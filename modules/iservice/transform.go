package iservice

import ("github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/errors"
)

func ExtractSvDefineTx(data []byte) (interface{}, error) {
	tx, err := sdk.LoadTx(data)
	if err != nil {
		return nil, err
	}
	txl, ok := tx.Unwrap().(sdk.TxLayer)
	for ok {
		tx = txl.Next()
		txl, ok = tx.Unwrap().(sdk.TxLayer)
	}
	ctx, ok := tx.Unwrap().(TxDefineService)
	if !ok {
		return nil, errors.ErrUnknownTxType(tx)
	}
	// now reformat this....
	return &ctx, nil
}