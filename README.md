# Go Flargs

<img src="go-flargs-gopher-again.png" alt="go flargs" title="go flargs" height="250" />

Flargs is a simple and lightweight framework for building command-line programs with the following design goals:

1. Is testable, providing abstractions around stdin, stdout, stderr, etc
2. Decouples the act of parsing arguments from the act of consuming inputs
3. Is chainable and composable, allowing for arbitrarily large and complex apps

Flargs conceives of 3 lifecycles, cleanly seperated:

1. *Parsing Flags and Args (flarging)*. This is the act of parsing arguments and flags into a custom structure (a flarg). The step allows no access to the environment.
2. *Loading flargs*. This step allows access to an environment is allows further processing and validating.
3. *Execution*. This is where your command is run. It runs against the object you created in step 1 and 2.


Flargs is composed of 3 basic components:

## Konf

This is your custom object which your app will run against. It takes any shape you want, but you must embed `flargs.StateMachine`:

```go
type myApiClient struct {
    hostname string
    port int
    path string
    flargs.stateMachine
}
```

Because you've embedded `flargs.StateMachine`, the struct will automatically implement this interface:

```go
type Flarger[T any] interface {
	Parse([]string) ([]string, error)
	Load(*Environment) error
    Run(*Environment) error
}
```
But you will want to define at least one of these on your own to get any interesting behaviour.

## Environment

An execution environment representing all the inputs and outputs a CLI should need.

```go
type Environment struct {
	InputStream  io.ReadWriter
	OutputStream io.ReadWriter
	ErrorStream  io.ReadWriter
	Randomness   io.Reader
	Filesystem   fs.FS
	Variables    map[string]string
}
```

This object is injected using dependency injection. Your CLI must use it for all i/o. So:

```go
// badly behaved ☹ don't do it
if os.Getenv("USER") == "sam" {
	fmt.Println("Sam, I am")
}

// well behaved ☺ this is the way
if env.Variables["USER"] == "sam" {
	fmt.Fprintln(env.OutputStream, "Sam, I am")
}
```

## Command

A Command is a Konf plus an Environment, along with a way to run the former against the latter. It has `Pipe()` for composability and a handful of helper methods.

```go
type Command[T any] struct {
	Env     *Environment
	runFunc RunFunc[T]
}

func (com1 Command[T]) Pipe(conf1 T, env2 *Environment) error {
	...
}
func (c Command) ParseAndLoad(args []string) error {
	...
}
```

# Getting Started

A simple hello-world program that allows you to swap "world" for something else might look like this:

```go
import (
    "github.com/sean9999/go-flargs"
)

//  our input structure. we only care about one value: name
type helloConf struct {
    name string
    flargs.StateMachine
}

//  get arg, set name
func (c *helloConf) Parse(args []string) error {
    if len(args) > 1 {
        return errors.New("too many args")
    }
    if len(args) == 1 {
        c.name = args[0]
    }
    c.name = "world"
    return nil
}

//  say hello
func (c *helloConf) Run(env *Environment) error {
    fmt.Fprintf(env.OutputStream, "hello %s", c.name)
}


conf := new(helloConf)
env := flargs.NewCLIEnvironment();
cmd := flargs.NewCommand(env, conf)
cmd.Run()
```

This might look pretty verbose for a simple CLI. But we now have a hermetic app that can be easily tested. It can grow in complexity without extra overhead. To test, we might do this:

```go
func TestNewCommand_hello(t *testing.T) {

    conf := new(helloConf)

    //  run command in testing mode
	env := flargs.NewTestingEnvironment(nil)
	cmd := flargs.NewCommand(env, conf)
    cmd.Parse([]string{"robin"})
	cmd.Run(conf)

    //  expected output
    want := "hello, robin"

    //  actual output
	got := env.GetOutput()

    if want != got.String() {
        t.Errorf("wanted %q but got %q", want, string(got))
    }
}
```
