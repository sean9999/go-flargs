package flargs

import "io"

// a Command is a Flarger with an Environment
type Command struct {
	Flarger
	*Environment
}

func NewCommand(fl Flarger, env *Environment) Command {
	return Command{fl, env}
}

func (k Command) ParseAndLoad(args []string) error {
	err := k.Parse(args)
	if err != nil {
		return err
	}
	err = k.Load()
	return err
}

func (k Command) LoadAndRun() error {
	err := k.Load()
	if err != nil {
		return err
	}
	return k.Run()
}

func (k Command) Load() error {
	return k.Flarger.Load(k.Environment)
}

func (k Command) Run() error {
	return k.Flarger.Run(k.Environment)
}

func Pipe(f1 Command, f2 Command) (int64, error) {
	f1.Run()
	return io.Copy(f2.Environment.InputStream, f1.Environment.OutputStream)
}

// // RunFunc is the core functionality, accepting T as input and operating against an [Environment]
// type RunFunc[T Flarger[T]] func(*Environment, T) error

// // a Command is a [RunFunc], along with an [Environment] to operate on.
// type Command[T Flarger[T]] struct {
// 	Env     *Environment
// 	runFunc RunFunc[T]
// }

// // Run runs the [RunFunc] against its *[Environment]
// func (com Command[T]) Run(conf T) error {
// 	return com.runFunc(com.Env, conf)
// }

// // Pipe is a convenience method for piping one Command to another.
// func (com Command[T]) Pipe(conf1 T, destEnvironment *Environment) error {
// 	err := com.Run(conf1)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = io.Copy(destEnvironment.InputStream, com.Env.OutputStream)
// 	if err != nil {

// 		return err
// 	}
// 	return nil
// }

// // NewCommand creates a new Command by combining an [Environment] and [RunFunc]
// func NewCommand[T Flarger[T]](env *Environment, runFn RunFunc[T]) Command[T] {
// 	com := Command[T]{
// 		Env:     env,
// 		runFunc: runFn,
// 	}
// 	return com
// }
