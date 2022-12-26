package utils

func InArray[T comparable](needle T, haystack []T) bool {
	key := needle

	for _, item := range haystack {
		if key == item {
			return true
		}
	}

	return false
}
