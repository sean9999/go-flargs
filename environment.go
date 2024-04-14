package flargs

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"strings"
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

// NewCLIEnvironment peoduces an Environment suitable for a CLI
func NewCLIEnvironment() *Environment {
	variables := map[string]string{
		"FLARGS_VERSION":         "v0.1.1",
		"FLARGS_EXE_ENVIRONMENT": "cli",
	}
	kvs := os.Environ()
	//	import environment
	for _, kv := range kvs {
		parts := strings.Split(kv, "=")
		if len(parts) == 2 {
			variables[string(parts[0])] = string(parts[1])
		}
	}

	env := Environment{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Randomness:   rand.Reader,
		Variables: map[string]string{
			"FLARGS_VERSION":         "v0.1.0",
			"FLARGS_EXE_ENVIRONMENT": "CLI",
		},
	}
	return &env
}

// NewTestingEnvironment produces an [Environment] suitable for testing.
// Pass in a "randomnessProvider" that offers a level of determinism that works for you.
// For good ole fashioned regular randomness, pass in [rand.Reader]
func NewTestingEnvironment(randomnessProvider io.Reader) *Environment {
	env := Environment{
		InputStream:  new(bytes.Buffer),
		OutputStream: new(bytes.Buffer),
		ErrorStream:  new(bytes.Buffer),
		Randomness:   randomnessProvider,
		Variables: map[string]string{
			"FLARGS_VERSION":         "v0.1.0",
			"FLARGS_EXE_ENVIRONMENT": "TESTING",
		},
	}
	return &env
}
