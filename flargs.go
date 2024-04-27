package flargs

type Flarger[T any] interface {
	Parse([]string) ([]string, error)
	Load(*Environment) error
}

var AnyFlarger Flarger[any]

// ParseFunc[T] takes a slice of command-line arguments.
// It returns a well-formed T, a "tail" of unused args, and an error
// It's run at parse-time
type ParseFunc[T Flarger[T]] func([]string) ([]string, error)

// LoadFunc does further processing on a Konfig.
// It is run at command-time, having access to an Environment
type LoadFunc[T Flarger[T]] func(*Environment) error

// flargMachine implements Flarger
type flargMachine[T Flarger[T]] struct {
	parseFunc ParseFunc[T]
	loadFunc  LoadFunc[T]
	konfig    T
}

func (m flargMachine[T]) Parse(args []string) ([]string, error) {
	return m.parseFunc(args)
}
func (m flargMachine[T]) Load(env *Environment) error {
	return m.loadFunc(env)
}
func (m flargMachine[T]) Flargs() T {
	return m.konfig
}

func NewFlarger[T Flarger[T]](pFunc ParseFunc[T], lFunc LoadFunc[T], konfig T) Flarger[T] {
	machine := flargMachine[T]{pFunc, lFunc, konfig}
	return machine
}
