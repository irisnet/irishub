package types

type TxCounter interface {
	Count() int64
	Incr()
}

type ValidTxCounter struct {
	count int64
}

func (vtc *ValidTxCounter) Count() int64 {
	return vtc.count
}

func (vtc *ValidTxCounter) Incr() {
	vtc.count++
}

func NewValidTxCounter() TxCounter {
	return &ValidTxCounter{
		count: 0,
	}
}
