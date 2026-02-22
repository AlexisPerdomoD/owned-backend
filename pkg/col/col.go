// Package col provides a collection of types and functions for collections.
package col

// Set is a generic collection of unique elements.
//
// It is implemented as a map[T]struct{}, where the zero-size struct is used
// to minimize memory overhead.
//
// Properties:
//   - Elements are unique (no duplicates).
//   - Average O(1) time complexity for Add, Remove, and Contains.
//   - Iteration order is not guaranteed.
//   - The set must be initialized before use (e.g. with NewSet).
//
// Zero value:
//   - The zero value of Set is nil and is NOT ready for use.
//   - Calling methods that mutate the set on a nil Set will panic.
type Set[T comparable] map[T]struct{}

// NewSet creates and returns an empty initialized Set.
//
// Example usage:
//
//	s := col.NewSet[int]()
func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

// Add inserts an item into the set.
//
// Returns true if the item was not already present.
// Returns false if the item already existed and the set was not modified.
//
// Complexity: O(1) average.
func (s Set[T]) Add(item T) bool {
	_, ok := s[item]
	if ok {
		return false
	}

	s[item] = struct{}{}
	return true
}

// Contains reports whether the item exists in the set.
//
// Complexity: O(1) average.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Remove deletes an item from the set.
//
// Returns true if the item existed and was removed.
// Returns false if the item was not present.
//
// Complexity: O(1) average.
func (s Set[T]) Remove(item T) bool {
	_, ok := s[item]
	if !ok {
		return false
	}

	delete(s, item)
	return true
}

// Slice returns all elements of the set as a slice.
//
// The order of elements in the returned slice is undefined.
//
// Complexity: O(n).
func (s Set[T]) Slice() []T {
	out := make([]T, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}
