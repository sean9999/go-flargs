package flargs

import (
	"io"
)

// RunFunc is the core functionality, accepting T as input and operating against an [Environment]
type RunFunc[T Flarger[T]] func(*Environment, T) error

// a Command is a [RunFunc], along with an [Environment] to operate on.
type Command[T Flarger[T]] struct {
	Env     *Environment
	runFunc RunFunc[T]
}

func (com Command[T]) Load(conf T) error {
	return conf.Load(com.Env)
}

// Run runs the [RunFunc] against its *[Environment]
func (com Command[T]) Run(conf T) error {
	return com.runFunc(com.Env, conf)
}

// Pipe is a convenience method for piping one Command to another.
func (com Command[T]) Pipe(conf1 T, destEnvironment *Environment) error {
	err := com.Run(conf1)
	if err != nil {
		return err
	}
	_, err = io.Copy(destEnvironment.InputStream, com.Env.OutputStream)
	if err != nil {

		return err
	}
	return nil
}

// NewCommand creates a new Command by combining an [Environment] and [RunFunc]
func NewCommand[T Flarger[T]](env *Environment, runFn RunFunc[T]) Command[T] {
	com := Command[T]{
		Env:     env,
		runFunc: runFn,
	}
	return com
}
