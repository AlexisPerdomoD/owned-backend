package helper_pkg

import "sync"

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

// MapConcurrentOutput represents the result of a concurrent mapping operation,
// containing both the computed value and any error that occurred during computation.
type MapConcurrentOutput[R any] struct {
	Value R
	Error error
}

// IsOk returns true if no error occurred during the concurrent mapping operation.
func (o *MapConcurrentOutput[R]) IsOk() bool {
	return o.Error == nil
}

// MapConcurrent applies the callback function cb concurrently to each element in the input slice
// and returns a slice of MapConcurrentOutput containing the results and any errors.
// The maxRoutines parameter controls the maximum number of concurrent goroutines.
// If maxRoutines is <= 0, a default limit of 10 concurrent routines is used.
// If the input slice is nil, it returns nil. If the input slice is empty,
// it returns an empty slice.
func MapConcurrent[T any, R any](input []T, cb func(T) (R, error), maxRoutines int) []MapConcurrentOutput[R] {
	if input == nil {
		return nil
	}

	out := make([]MapConcurrentOutput[R], len(input))
	if len(input) == 0 {
		return out
	}

	if maxRoutines <= 0 {
		// this is a dumpt safe limit that does not care for the actual concurrent task
		// to work in concurrent, please make sure limit is place accordingly to the tasks
		maxRoutines = 10
	}

	wg := sync.WaitGroup{}
	limitter := make(chan struct{}, maxRoutines)

	for i, v := range input {
		wg.Add(1)
		limitter <- struct{}{}

		go func(i int, v T) {
			defer wg.Done()
			val, err := cb(v)
			out[i].Value = val
			out[i].Error = err
			<-limitter
		}(i, v)
	}

	wg.Wait()
	return out
}
