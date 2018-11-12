package keys

import "fmt"

func ErrKeyNameConflict(name string) error {
	return fmt.Errorf("acount with name %s already exists", name)
}

func ErrMissingName() error {
	return fmt.Errorf("you have to specify a name for the locally stored account")
}

func ErrMissingPassword() error {
	return fmt.Errorf("you have to specify a password for the locally stored account")
}

func ErrMissingSeed() error {
	return fmt.Errorf("you have to specify seed for key recover")
}