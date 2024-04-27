package flargs

type Flarger[T any] interface {
	Parse([]string) ([]string, error)
	Load(*Environment) error
}
