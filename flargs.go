package flargs

// ParseFunc[T] takes a slice of command-line arguments.
// It returns a well-formed T, a "tail" of unused args, and an error
type ParseFunc[T any] func([]string) (T, []string, error)
