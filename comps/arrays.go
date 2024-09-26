package comps

// UnorderedEquals checks if two arrays contain the same elements, regardless of the order.
func UnorderedEquals[T comparable](arr1, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	counts := make(map[T]int)
	for _, item := range arr1 {
		counts[item]++
	}

	for _, item := range arr2 {
		if counts[item] == 0 {
			return false
		}
		counts[item]--
	}

	return true
}
