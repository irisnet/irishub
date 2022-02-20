package types

const (
	DoNotModify = "[do-not-modify]"
)

// Modified returns whether the field is modified
func Modified(target string) bool {
	return target != DoNotModify
}