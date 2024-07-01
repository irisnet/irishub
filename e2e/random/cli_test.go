package random

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}