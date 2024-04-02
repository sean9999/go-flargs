package flargs

type RunFunc[T any] func(*Environment, T)

type Command[T any] struct {
	Env       *Environment
	ParseFunc ParseFunc[T]
	RunFunc   RunFunc[T]
}

func (com Command[T]) Run(rawArgs []string) ([]string, error) {
	niceArgs, tail, err := com.ParseFunc(rawArgs)
	com.RunFunc(com.Env, niceArgs)
	return tail, err
}

func NewCommand[T any](env *Environment, parseFn ParseFunc[T], runFn RunFunc[T]) Command[T] {
	com := Command[T]{
		Env:       env,
		ParseFunc: parseFn,
		RunFunc:   runFn,
	}
	return com
}
