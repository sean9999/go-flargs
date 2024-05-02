package main

import (
	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/rot13"
)

func main() {

	params := new(rot13.RotKonf)
	env := flargs.NewCLIEnvironment("")
	cmd := flargs.NewCommand(params, env)
	err := cmd.ParseAndLoad(nil)
	if err != nil {
		panic(err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

}
