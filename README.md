# Go Flargs

<img src="go-flargs-gopher-again.png" alt="go flargs" title="go flargs" height="250" />

Flargs is a simple and lightweight framework for building command-line programs with the following design goals:

1. Is testable, providing abstractions around stdin, stdout, stderr, etc
2. Decouples the act of parsing arguments from the act of consuming inputs
3. Is chainable and composable, allowing for arbitrarily large and complex apps

Flargs conceives of n lifecycles, cleanly seperated:

1. *Parsing Flags and Args (flarging)*. This is the act of parsing arguments and flags into a custom structure (a flarg). The step allows no access to the environment.
2. *Loading flargs*. This step allows access to an environment
3. *Execution*. This is where your command is run.


Flargs is composed of 4 basic components:

## The Flarger Interface

This is your custom object. You can decide what it looks like, but it must satisfy this interface:

```go
type Flarger[T any] interface {
	Parse([]string) ([]string, error)
	Load(*Environment) error
}
```

## ParseFunc

A ParseFunc takes in a slice of strings and produces a structure that you define. It consumes some flags and args, and leaves others unparsed, returning them for later parsing. It can return an error indicating the flags and args were insufficient to run your Command.

Its signature is:

```go
type ParseFunc[T any] func([]string) (T, []string, error)
```

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

## RunFunc

This is the meat of the functionality and its signature is:

```go
type RunFunc[T any] func(*Environment, T) error
```

a RunFunc should be hermetic. It should read from `Environment.InputStream` and write to `Environment.Outputsteam`. Although it returns an error, any error information meant to be displayed on a terminal should be sent to `Environment.ErrorStream`.

If you violate these principles, you won't have a happy time. You will not be able to take advantage of the true power of Flargs.

```go
// badly behaved ☹
if os.Getenv("USER") == "sam" {
	fmt.Println("Sam, I am")
}

// well behaved ☺
if env.Variables["USER"] == "sam" {
	fmt.Fprintln(env.OutputStream, "Sam, I am")
}
```

## Command

A Command is a RunFunc plus an Environment, along with a way to run the former against the latter. It also has a `Pipe()` for composability.

```go
type Command[T any] struct {
	Env     *Environment
	runFunc RunFunc[T]
}

func (com Command[T]) Run(conf T) error {
	...
}
func (com1 Command[T]) Pipe(conf1 T, env2 *Environment) error {
	...
}
```

# Getting Started

A simple hello-world program that allows you to swap "world" for something else might look like this:

```go
import (
    "github.com/sean9999/go-flargs"
)

//  our input structure. we only care about one flag.
type conf struct {
    name string
}

//  this is our ParseFunc. It returns a *conf
parseFn := func(args []string) (*conf, []string, error) {
    params := new(conf)
    //  default value
    params.name = "world"
    fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
    fset.Func("name", "hello to who?", func(s string) error {
        if s == "batman" {
            return errors.New("you cannot say hello to batman")
        }
		params.name = s
		return nil
	})
    err := fset.Parse()
    return params, fset.Args(), err
}

// this is our RunFunc. it says hello to params.name.
// it writes to env.OutputStream, which in a real CLI is os.Stdout
helloFn := func(env *flargs.Environment, params *conf) error {
    fmt.Fprintf(env.OutputStream, "hello, %s", params.name)
    return nil
}

//  first parse
conf, _, err := parseFn(os.Args())
if err != nil {
    panic(err)
}

//  if all goes well, run
env := flargs.NewCLIEnvironment();
cmd := flargs.NewCommand(env, helloFn)
cmd.Run(conf)
```

This might look pretty verbose for a simple CLI. But we now have a hermetic app that can be easily tested. It can grow in complexity without extra overhead. To test, we might do this:

```go
func TestNewCommand_hello(t *testing.T) {

    //  parse args
	conf, _, _ := parseFn([]string{"--name", "robin"})

    //  run command in testing mode
	env := flargs.NewTestingEnvironment(nil)
	cmd := flargs.NewCommand(env, helloFn)
	cmd.Run(conf)

    //  expected output
    want := "hello, robin"

    //  actual output
	got := cmd.Env.GetOutput()

    if want != got.String() {
        t.Errorf("wanted %q but got %q", want, string(got))
    }
}
```
