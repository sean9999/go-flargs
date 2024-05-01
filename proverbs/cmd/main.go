package main

import (
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
)

func main() {

	params := new(proverbs.Params)
	env := flargs.NewCLIEnvironment("")
	cmd := flargs.NewCommand(params, env)
	err := cmd.ParseAndLoad(os.Args[1:])
	if err != nil {
		panic(err)
	}
	cmd.Run()

}
