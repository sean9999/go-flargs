package main

import (
	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/noop"
)

func main() {

	//	noop is a command that does absolutely nothing
	konf := new(noop.NoopConf)
	cmd := flargs.NewCommand(konf, nil)
	err := cmd.ParseAndLoad([]string{})
	if err != nil {
		panic(err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

}
