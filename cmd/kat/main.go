package main

import (
	"io/fs"
	"os"

	"github.com/sean9999/go-flargs"
)

// KatConf is all the information necessary to parse flags for, and then prepare to execute kat
type KatConf struct {
	fileNames     []string  // filled at parse-time
	files         []fs.File // created at command-time
	withNumbering bool
	remainingArgs []string
}

func (k *KatConf) Args() []string {
	return k.remainingArgs
}

func main() {

	konf := new(KatConf)
	env := flargs.NewCLIEnvironment("/")
	katCmd := flargs.NewCommand(konf, env)

	err := katCmd.ParseAndLoad(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	err = katCmd.LoadAndRun()
	if err != nil {
		os.Exit(2)
	}

}
