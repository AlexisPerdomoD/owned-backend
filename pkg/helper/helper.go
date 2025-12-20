package helper

import "fmt"

// Map applies the callback function cb to each element in the input slice and returns a new slice
// containing the results. If the input slice is nil, it returns nil. If the input slice is empty,
// it returns an empty slice.
func Map[T any, R any](input []T, cb func(T) R) []R {
	if input == nil {
		return nil
	}

	out := make([]R, len(input))
	if len(input) == 0 {
		return out
	}

	for i, v := range input {
		out[i] = cb(v)
	}

	return out
}

func AssertNotNil(v any, entity string) {
	if any(v) == nil {
		panic(fmt.Sprintf("%s provided as nil", entity))
	}
}
