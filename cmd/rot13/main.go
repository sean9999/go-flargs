package main

import (
	"io/fs"
	"os"

	"github.com/joshlf13/rot13"
	"github.com/sean9999/go-flargs"
)

type state struct {
	fileName      string
	inputText     []byte
	remainingArgs []string
}

func (s *state) Args() []string {
	return s.remainingArgs
}
func (s *state) Parse(args []string) error {
	//	let's be naive and take any input
	if len(args) > 0 {
		s.fileName = args[0]
		s.remainingArgs = args[1:]
	}
	return nil
}
func (s *state) Load(env *flargs.Environment) error {
	if s.fileName != "" {
		contents, err := fs.ReadFile(env.Filesystem, s.fileName)
		if err != nil {
			return err
		}
		s.inputText = contents
		return nil
	} else {
		s.inputText = env.GetInput()
		return nil
	}
}
func (s *state) Run(env *flargs.Environment) error {

	wr := rot13.NewWriter(env.OutputStream)
	_, err := wr.Write(s.inputText)
	return err

}

func main() {

	params := new(state)
	env := flargs.NewCLIEnvironment(".")
	cmd := flargs.NewCommand(params, env)
	err := cmd.ParseAndLoad(os.Args[1:])
	if err != nil {
		panic(err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

}
