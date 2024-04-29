package rot13

import (
	"io/fs"

	"github.com/joshlf13/rot13"
	"github.com/sean9999/go-flargs"
)

type State struct {
	fileName      string
	inputText     []byte
	remainingArgs []string
}

func (s *State) Args() []string {
	return s.remainingArgs
}
func (s *State) Parse(args []string) error {
	//	let's be naive and take any input
	if len(args) > 0 {
		s.fileName = args[0]
		s.remainingArgs = args[1:]
	}
	return nil
}
func (s *State) Load(env *flargs.Environment) error {
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
func (s *State) Run(env *flargs.Environment) error {
	wr := rot13.NewWriter(env.OutputStream)
	_, err := wr.Write(s.inputText)
	return err
}
