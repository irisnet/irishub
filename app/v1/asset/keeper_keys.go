package asset

var (
	// the separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes

	// TODO
)

func getStringsKey(ss []string) (result []byte) {
	for _, s := range ss {
		result = append(append(
			result,
			[]byte(s)...),
			emptyByte...)
	}
	if len(result) > 0 {
		return result[0 : len(result)-1]
	}
	return
}
