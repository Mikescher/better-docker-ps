package langext

func Ptr[T any](v T) *T {
	return &v
}

func Coalesce[T any](a *T, b T) T {
	if a != nil {
		return *a
	}
	return b
}
