package slices

import "reflect"

// SlicesAreEquivalent checks if two slices of strings are equivalent,
// treating nil and empty slices as equivalent.
func SlicesAreEquivalent(a, b []string) bool {
	return (len(a) == 0 && len(b) == 0) || reflect.DeepEqual(a, b)
}
