# Go Flargs

<img src="go-flargs-gopher-again.png" alt="go flargs" title="go flargs" height="250" />

Flargs is an opinionated package for building command-line programs with the following design goals:

1. Is testable, providing abstractions around stdin, stdout, stderr, etc
2. Removes complexity of argument parsing
3. Decouples the act of parsing arguments from the act of consuming inputs
4. Provides a nice, sane, clean interface
5. Is chainable and composable, allowing for arbitrarily large and complex apps

Flargs conceives of two lifecycles, cleanly seperated:

1. Parsing flags and args into a sensible structure. You get to define what this looks like.
2. Running the code. This lifecycle knows nothing about args and flags. It only knows about the sensible structure that was passed to it.

You can choose to throw errors early on the parsing stage, if the args and flags don't make sense, or later on in the command execution phase. 

Flargs is composed of 4 basic components:

## ParseFunc

This is a function taking in a slice of strings (such as `os.Args()`) and producing an object that makes sense for your command (using generics), along with an error, and unparsed arguments (enabling composibility). Its signature is:

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
	Variables    map[string]string // environment variables
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

type conf struct {
    name string
}

//  this is a flargs.ParseFunc
parseFn := func(args []string) (*conf, []string, error) {
    conf := new(conf)
    //  default value
    conf.name = "world"
    fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
    fset.Func("name", "hello to who?", func(s string) error {
        if s == "batman" {
            return errors.New("you cannot say hello to batman")
        }
		conf.name = s
		return nil
	})
    err := fset.Parse()
    return conf, fset.Args(), err
}

//  this is a flargs.RunFunc
helloFn := func(env *flargs.Environment, conf *catConf) error {
    outputString := fmt.Sprintf("hello, %s")
    env.OutputStream.Write([]byte(outputString))
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

This might look pretty verbose for a simple CLI. But we now have a hermetic app that can be easily tested. It can grow in complexity without extra overhead. We've added all the _necessary_ complexity already. To test, we might do this:

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
	got := new(bytes.Buffer)
	got.ReadFrom(cmd.Env.OutputStream)

    if want != got.String() {
        t.Errorf("wanted %q but got %q", want, got.String())
    }
}
```
