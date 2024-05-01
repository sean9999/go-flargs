package hello

import (
	"errors"
	"fmt"

	"github.com/sean9999/go-flargs"
)

type params struct {
	name string
}

type Konf struct {
	params
	flargs.StateMachine
}

func (c Konf) Parse(args []string) error {
	if len(args) == 1 {
		c.name = args[0]
		c.RemainingArgs = args[1:]
		c.Phase = flargs.Parsing
		return nil
	}
	return errors.New("wrong number of args")
}

func (c Konf) Run(env *flargs.Environment) error {
	fmt.Fprintf(env.OutputStream, "hello, %#v", c)
	return nil
}
