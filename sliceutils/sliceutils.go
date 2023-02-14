package sliceutils //TODO: this will soon be replaced by the built in slices package

func RemoveByIndex[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
