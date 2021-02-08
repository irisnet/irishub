package types

// HasDuplicate checks if the given array contains duplicate elements.
// Return true if there exist duplicate elements, false otherwise
func HasDuplicate(arr []string) bool {
	elementMap := make(map[string]bool)

	for _, elem := range arr {
		if _, ok := elementMap[elem]; ok {
			return true
		}

		elementMap[elem] = true
	}

	return false
}
