package flargs

import (
	"io"
)

// RunFunc is the core functionality, accepting T as input
type RunFunc[T any] func(*Environment, T) error

// a Command is a [RunFunc], along with an [Environment] to operate on.
type Command[T any] struct {
	Env     *Environment
	runFunc RunFunc[T]
}

// Run runs the [RunFunc] against its *[Environment]
func (com Command[T]) Run(conf T) error {
	return com.runFunc(com.Env, conf)
}

// Pipe is a convenience method for piping one Command to another
func (com1 Command[T]) Pipe(conf1 T, env2 *Environment) error {
	err := com1.Run(conf1)
	if err != nil {
		return err
	}
	_, err = io.Copy(env2.InputStream, com1.Env.OutputStream)
	if err != nil {
		return err
	}
	return nil
}

// NewCommand creates a new Command by combining an [Environment] and [RunFunc]
func NewCommand[T any](env *Environment, runFn RunFunc[T]) Command[T] {
	com := Command[T]{
		Env:     env,
		runFunc: runFn,
	}
	return com
}
