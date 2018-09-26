package iservice

import "testing"

type TESTEnumA int
type TESTEnumB int

const (
	AOption TESTEnumA = 1
	BOption TESTEnumA = 2
	COption TESTEnumB = 3
	DOption TESTEnumB = 4
)


func TestEnum(t *testing.T) {

	println(AOption)
	println(BOption)
	println(COption)
	println(DOption)
}