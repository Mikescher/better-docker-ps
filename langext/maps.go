package langext

func MapKeyArr[T comparable, V any](v map[T]V) []T {
	result := make([]T, 0, len(v))
	for k := range v {
		result = append(result, k)
	}
	return result
}
