package flargs

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
)

// Enviroment is an execution environment for a Command.
// In the context of a CLI, these would be [os.StdIn], [os.StdOut], etc
// In the context of a test-suite, they would probably be [bytes.Buffer]
type Environment struct {
	InputStream  io.ReadWriter
	OutputStream io.ReadWriter
	ErrorStream  io.ReadWriter
	Randomness   io.Reader
	Variables    map[string]string
}

func NewCLIEnvironment() *Environment {
	env := Environment{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Randomness:   rand.Reader,
		Variables: map[string]string{
			"PLATOON_VERSION":         "v0.1.0",
			"PLATOON_EXE_ENVIRONMENT": "CLI",
		},
	}
	return &env
}

func NewTestingEnvironment() *Environment {
	env := Environment{
		InputStream:  new(bytes.Buffer),
		OutputStream: new(bytes.Buffer),
		ErrorStream:  new(bytes.Buffer),
		Randomness:   rand.Reader,
		Variables: map[string]string{
			"PLATOON_VERSION":         "v0.1.0",
			"PLATOON_EXE_ENVIRONMENT": "TESTING",
		},
	}
	return &env
}
