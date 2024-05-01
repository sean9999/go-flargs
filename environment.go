package flargs

import (
	"bytes"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"strings"
	"testing/fstest"
	"time"

	realfs "github.com/sean9999/go-real-fs"
)

// Enviroment is an execution environment for a Command.
// In the context of a CLI, these would be [os.StdIn], [os.StdOut], etc.
// In the context of a test-suite, you can use [bytes.Buffer] and [fstest.MapFS].
// For benchmarking, you can use a [NullDevice].
type Environment struct {
	InputStream  io.ReadWriter
	OutputStream io.ReadWriter
	ErrorStream  io.ReadWriter
	Randomness   rand.Source
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
		Randomness:   rand.NewSource(time.Now().UnixNano()),
		Filesystem:   realFs,
		Variables:    vars,
	}
	return &env
}

// NewTestingEnvironment produces an [Environment] suitable for testing.
// Pass in a "randomnessProvider" that offers a level of determinism that works for you.
// For good ole fashioned regular randomness, pass in [rand.Reader]
// If your program doesn't use randomness, just pass in nil.
func NewTestingEnvironment(randomnessProvider rand.Source) *Environment {
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

type NullDevice struct {
	io.Writer
}

func (b NullDevice) Read(_ []byte) (int, error) {
	return 0, nil
}
func (b NullDevice) Open(_ string) (fs.File, error) {
	return nil, nil
}

func NewNullEnvironment() *Environment {
	e := Environment{
		InputStream:  NullDevice{io.Discard},
		OutputStream: NullDevice{io.Discard},
		ErrorStream:  NullDevice{io.Discard},
		Filesystem:   NullDevice{},
		Variables:    map[string]string{},
	}
	return &e
}
