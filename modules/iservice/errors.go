// nolint
package iservice

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/errors"
)

var (
	errServiceNameEmpty = fmt.Errorf("Invalid service name")
	errServiceDescEmpty = fmt.Errorf("Invalid service description")
	errServiceExists    = fmt.Errorf("Service already exists")

	invalidInput = errors.CodeTypeBaseInvalidInput
)

func ErrServiceExists() error {
	return errors.WithCode(errServiceExists, errors.CodeTypeBaseInvalidInput)
}
func ErrMissingSignature() error {
	return errors.ErrMissingSignature()
}
