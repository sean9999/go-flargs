package flargs

import (
	"io"
)

// RunFunc executes code against streams defined in *Environment
type RunFunc[T any] func(*Environment, T) error

// Command is a container for a RunFunc and and *Environment
type Command[T any] struct {
	Env     *Environment
	RunFunc RunFunc[T]
}

// Run runs the RunFunc against the *Environment
func (com Command[T]) Run(conf T) error {
	return com.RunFunc(com.Env, conf)
}

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

// NewCommand creates a new command, taking in T as input
func NewCommand[T any](env *Environment, runFn RunFunc[T]) Command[T] {
	com := Command[T]{
		Env:     env,
		RunFunc: runFn,
	}
	return com
}
