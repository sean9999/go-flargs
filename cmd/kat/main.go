package main

import (
	"os"

	"github.com/sean9999/go-flargs"
)

func main() {

	konf := new(KatConf)

	_, err := konf.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	//	run command
	env := flargs.NewCLIEnvironment("/")
	konf.Load(env)
	cmd := flargs.NewCommand(env, KatFunction)
	err = cmd.Run(konf)

	//	exit code
	if err != nil {
		os.Exit(1)
	}

}
