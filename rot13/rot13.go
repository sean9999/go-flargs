package rot13

import (
	"bufio"

	"github.com/sean9999/go-flargs"
)

type State struct {
	fileName  string
	inputText []byte
	flargs.StateMachine
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
	runeStream := bufio.NewReader(env.InputStream)
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
