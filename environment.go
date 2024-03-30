package platoon

import "io"

// Enviroment is an execution environment for a Command.
// In the context of a CLI, these would be [os.StdIn], [os.StdOut], etc
// In the context of a test-suite, they would probably be [bytes.Buffer]
type Environment struct {
	InputStream  io.Reader
	OutputStream io.ReadWriter
	ErrorStream  io.ReadWriter
	Variables    map[string]string
}
