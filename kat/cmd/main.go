package main

import (
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/kat"
)

func main() {

	konf := new(kat.Konf)
	//pwd, _ := os.Getwd()
	env := flargs.NewCLIEnvironment(".")
	katCmd := flargs.NewCommand(konf, env)

	err := katCmd.ParseAndLoad(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	err = katCmd.Run()
	if err != nil {
		os.Exit(2)
	}

}
