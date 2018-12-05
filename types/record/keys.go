package record


import "fmt"

func KeyRecord(dataHash string) []byte {
	return []byte(fmt.Sprintf("record:%s", dataHash))
}

