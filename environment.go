package flargs

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	realfs "github.com/sean9999/go-real-fs"
)

// Enviroment is an execution environment for a Command.
// In the context of a CLI, these would be [os.StdIn], [os.StdOut], etc
// In the context of a test-suite, they would probably be [bytes.Buffer]
type Environment struct {
	InputStream  io.ReadWriter
	OutputStream io.ReadWriter
	ErrorStream  io.ReadWriter
	Randomness   io.Reader
	Filesystem   fs.FS
	Variables    map[string]string
}

func (e Environment) GetOutput() []byte {
	buf, _ := io.ReadAll(e.OutputStream)
	return buf
}

func (e Environment) GetError() []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(e.ErrorStream)
	return buf.Bytes()
}

func (e Environment) GetInput() []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(e.InputStream)
	return buf.Bytes()
}

// NewCLIEnvironment produces an Environment suitable for a CLI.
// It's a helper function with sane defaults.
func NewCLIEnvironment(baseDir string) *Environment {
	envAsMap := func(envs []string) map[string]string {
		m := make(map[string]string)
		i := 0
		for _, s := range envs {
			i = strings.IndexByte(s, '=')
			m[s[0:i]] = s[i+1:]
		}
		return m
	}

	//	import parent env vars
	vars := envAsMap(os.Environ())
	vars["FLARGS_VERSION"] = "v1.0.1"
	vars["FLARGS_EXE_ENVIRONMENT"] = "cli"

	realFs := realfs.New()

	env := Environment{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Randomness:   rand.Reader,
		Filesystem:   realFs,
		Variables:    vars,
	}
	return &env
}

// NewTestingEnvironment produces an [Environment] suitable for testing.
// Pass in a "randomnessProvider" that offers a level of determinism that works for you.
// For good ole fashioned regular randomness, pass in [rand.Reader]
// If your program doesn't use randomness, just pass in nil.
func NewTestingEnvironment(randomnessProvider io.Reader) *Environment {
	env := Environment{
		InputStream:  new(bytes.Buffer),
		OutputStream: new(bytes.Buffer),
		ErrorStream:  new(bytes.Buffer),
		Randomness:   randomnessProvider,
		Filesystem:   fstest.MapFS{},
		Variables: map[string]string{
			"FLARGS_VERSION":         "v1.0.1",
			"FLARGS_EXE_ENVIRONMENT": "testing",
		},
	}
	return &env
}
