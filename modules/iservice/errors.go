// nolint
package iservice

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/errors"
)

var (
	errServiceNameEmpty = fmt.Errorf("Invalid service name")
	errServiceDescEmpty = fmt.Errorf("Invalid service description")
	errServiceChainID = fmt.Errorf("Invalid service chainID")
	errServiceMessaging = fmt.Errorf("Invalid service messaging,enum:'Unicast/Multicast'")
	errServiceMethods = fmt.Errorf("Methods is empty")
	errServiceMethodID = fmt.Errorf("Invalid service methodID")
	errServiceMethodIDNotUnique = fmt.Errorf("methodID is not unique")
	errServiceMethodName = fmt.Errorf("Invalid methodName")
	errServiceMethodNameNotUnique = fmt.Errorf("methodName is not unique")
	errServiceMethodDescription = fmt.Errorf("Invalid service methodDescription")
	errServiceMethodInput = fmt.Errorf("Invalid service methodInput")
	errServiceMethodOutput = fmt.Errorf("Invalid service methodOutput")
	errServiceMethodOutputPrivacy = fmt.Errorf("Invalid service methodOutputPrivacy")
	errServiceExists    = fmt.Errorf("Service already exists")

	invalidInput = errors.CodeTypeBaseInvalidInput
)

func ErrServiceExists() error {
	return errors.WithCode(errServiceExists, errors.CodeTypeBaseInvalidInput)
}
func ErrMissingSignature() error {
	return errors.ErrMissingSignature()
}
