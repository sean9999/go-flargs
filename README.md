# Go Flargs

Go Flargs is a package for parsing command-line flags and arguments, and then running commands. It has the following design-goals:

1. Is testable
2. Removes the complexity of argument parsing
3. Provides a nice, sane, clean interface
4. Is chainable, allowing sub-commands and sub-sub-commands

It's composed of 3 basic components:

## Margs

Margs is simply a map representing how inputs to our program should look. Margs are simple key-value pairs whose definition looks like this:

```go
type margs = map[string]any
```

So that a sensible input if you were building git might look like this:

```go
margs := map[string]any{
    "subCommand": "fetch",
    "repo": "/some/folder",
    "branch": "main",
    "remote": "origin",
    "verbose:" false
}
```

## Flargs

The Flargs interface houses and calls a FlargParser, which is a function that converts `[]string` to a `marg`.

It also returns a tail (`[]string`), reprenting those arguments that this parser didn't care about. This can be used to allow one command to pass-off execution to another (ie: chaining).


## Commands

A Command takes in fully-formed Margs, and has access to en Environment which it uses to write to and read from. A Command does not know or care about arguments in `[]string` form. A well-behaved Command will not write to or read from anything outside its Environment.

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


## Environment

An execution environment representing all the inputs and outputs a CLI should have.
Testing environments will have different inputs and outputs than the CLI proper.

# Getting Started

A simple hello-world program might look like this:

```go
import (
    "github.com/sean9999/go-flargs"
)

//  a function that parses arguments into a map
parseFn := func(args []string) (map[string]any, []string, error) {
    fset := *flag.NewFlagSe("hello world", flag.ContinueOnError)
    //  set default values
    m := map[string]any{
        "name": "world",
    }
    fset.Func("name", "the name of the person to greet", func(s string) error {
        //  passing a blank name will be considered an error
        if len(s) == 0 {
            return errors.New("zero length name")
        }
        m["name"] = s
        return nil
    })
    err := fset.Parse(args)
    return m, nil, err
}
helloParser := flargs.NewFlargs()

//  flargs.Parse() gives us nicely formed Margs
margs, _, _ := helloParser.Parse([]string{"--name=bert"})

//  An execution environment
env := flargs.NewCLIEnvironment()

//  our hello world command, along with an execution environment
helloCmd := flargs.NewCommand(env, func(env *Environment, margs map[string]any) error {
    output := fmt.Sprintf("hello %s", margs["name"])
    
    //  this is how commands must output
    env.OutputSteam.Write(output)
    return nil
})

helloCmd.Run(margs)
```

This might look pretty verbose for a simple CLI. But we have gained some extra power and flexibility. We have massaged our inputs in a way that's more powerful than `flag.Parse()` alone, although we are free to use `flag.Parse()` at will. The inputs to our program take a more natural shape. And testing now is simple:

```go
import (
	"testing"

	"github.com/sean9999/go-platoon"
)

func TestHelloWorld(t *testing.T) {

    //  all the same code, except...

    //  *testing* execution environment
    env := platoon.NewTestingEnvironment()

    //  will write "hello bert" to env.OutputStream
    helloCmd.Run(margs)

    want := "hello bert"

    //  read env.OutputStream 
	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(env.OutputStream)
    got := buf.String()

    //  compare
    if want != got {
        t.Errorf("got %q but wanted %q", got, want)
    }

}
```
