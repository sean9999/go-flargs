package main

import (
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/hello"
)

func main() {

	//	hello is a command that says hello.
	//	You need to pass in one arg
	k := new(hello.Konf)
	err := k.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	env := flargs.NewCLIEnvironment("/")
	err = k.Run(env)

	if err != nil {
		panic(err)
	}

}
