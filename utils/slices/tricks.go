package slices

import "github.com/katallaxie/pkg/utils"

// Cut removes an element from a slice at a given position.
func Cut[T comparable](i, j int, a ...T) []T {
	copy(a[i:], a[j:])
	for k, n := len(a)-j+i, len(a); k < n; k++ {
		a[k] = utils.Zero[T]()
	}

	return a[:len(a)-j+i]
}

// Delete removes an element from a slice by value.
func Delete[T comparable](i int, a ...T) []T {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = utils.Zero[T]()

	return a[:len(a)-1]
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

// Filter removes all elements from a slice that satisfy a predicate.
func Filter[T comparable](f func(T) bool, a ...T) []T {
	b := a[:0]
	for _, x := range a {
		if !f(x) {
			b = append(b, x)
		}
	}

	return b
}
