package types

type ValidTxCounter struct {
	count int64
}

func (vtc *ValidTxCounter) Count() int64 {
	return vtc.count
}

func (vtc *ValidTxCounter) Incr() {
	vtc.count++
}

func NewValidTxCounter() *ValidTxCounter {
	return &ValidTxCounter{
		count: 0,
	}
}
