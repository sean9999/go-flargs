package main

import (
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/rot13"
)

func main() {

	params := new(rot13.State)
	env := flargs.NewCLIEnvironment(".")
	cmd := flargs.NewCommand(params, env)
	err := cmd.ParseAndLoad(os.Args[1:])
	if err != nil {
		panic(err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

}
