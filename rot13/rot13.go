package rot13

import (
	"bufio"
	"bytes"
	"io/fs"

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

	//	rotate a rune
	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (((r - 'A') + 13) % 26)
		case r >= 'a' && r <= 'z':
			return 'a' + (((r - 'a') + 13) % 26)
		default:
			return r
		}
	}

	//	read in input rune by rune
	runeStream := bufio.NewReader(bytes.NewReader(s.inputText))
	result := []byte{}
	for {
		if c, _, err := runeStream.ReadRune(); err != nil {
			break
		} else {
			result = append(result, byte(rot13(c)))
		}
	}

	//	write result
	env.OutputStream.Write(result)
	return nil

}
