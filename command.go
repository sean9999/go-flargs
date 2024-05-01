package flargs

import (
	"io"
)

// a Command is a Flarger with an [Environment]
type Command struct {
	Flarger
	*Environment
}

// creates a [Command]
func NewCommand(fl Flarger, env *Environment) Command {
	return Command{fl, env}
}

// combines [Command.Parse] and [Command.Load]
func (k Command) ParseAndLoad(args []string) error {
	err := k.Parse(args)
	if err != nil {
		return err
	}
	err = k.Load()
	return err
}

// LoadAndRun combines [Command.Load] and [Command.Run]
func (k Command) LoadAndRun() error {
	err := k.Load()
	if err != nil {
		return err
	}

	return k.Run()
}

// Load processes the flarg configuration in the context of an [Environment]
func (k Command) Load() error {
	return k.Flarger.Load(k.Environment)
}

// Run runs Flarger.Run in the context of an [Environment]
func (k Command) Run() error {
	return k.Flarger.Run(k.Environment)
}

// Pipe pipes one Command to another
func Pipe(f1 Command, f2 Command) (int64, error) {
	f1.Run()
	return io.Copy(f2.Environment.InputStream, f1.Environment.OutputStream)
}
