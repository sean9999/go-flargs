package flargs

// FlargsParseFunc[T] takes a slice of command-line arguments, returns a well-formed T
// a "tail" of unparsed args, and an error
type FlargsParseFunc[T any] func([]string) (T, []string, error)
