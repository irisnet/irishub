// reference: https://github.com/binance-chain/go-sdk/blob/master/common/common.go
package keystore

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	isAlphaNumFunc = regexp.MustCompile(`^[[:alnum:]]+$`).MatchString
)

func QueryParamToMap(qp interface{}) (map[string]string, error) {
	queryMap := make(map[string]string, 0)
	bz, err := json.Marshal(qp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bz, &queryMap)
	if err != nil {
		return nil, err
	}
	return queryMap, nil
}

func CombineSymbol(baseAssetSymbol, quoteAssetSymbol string) string {
	return fmt.Sprintf("%s_%s", baseAssetSymbol, quoteAssetSymbol)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func IsAlphaNum(s string) bool {
	return isAlphaNumFunc(s)
}
