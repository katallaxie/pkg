package cast

// Slice is a type that represents a slice.
func Slice[T any](val T) []T {
	return []T{val}
}
