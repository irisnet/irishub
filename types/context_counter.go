package types

type TxCounter interface {
	Count() int64
	Add()
}

type ValidTxCounter struct {
	count int64
}

func (vtc *ValidTxCounter) Count() int64 {
	return vtc.count
}

func (vtc *ValidTxCounter) Add() {
	vtc.count++
}

func NewValidTxCounter() TxCounter {
	return &ValidTxCounter{
		count: 0,
	}
}
