package slices

// Cut removes an element from a slice at a given position.
func Cut[T comparable](start, end int, a ...T) []T {
	return append(a[:start], a[end:]...)
}

// Delete removes an element from a slice by value.
func Delete[T comparable](idx int, a ...T) []T {
	return append(a[:idx], a[idx+1:]...)
}

// Push adds an element to the end of a slice.
func Push[T comparable](x T, a ...T) []T {
	return append(a, x)
}

// Pop removes an element from the end of a slice.
func Pop[T comparable](a ...T) (T, []T) {
	return a[len(a)-1], a[:len(a)-1]
}

// Insert adds an element at a given position in a slice.
func Insert[T comparable](x T, idx int, a ...T) []T {
	return append(a[:idx], append([]T{x}, a[idx:]...)...)
}
