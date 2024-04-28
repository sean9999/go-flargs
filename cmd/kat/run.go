package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sean9999/go-flargs"
)

// Now we should have an array of real files.
// We can iterate over them and spit out the contents.
// We still report an error if there is a problem operating on these real files.
func (k *KatConf) Run(env *flargs.Environment) error {
	var lastKnownError error
	line := 0
	//	stdin
	if k.withNumbering {
		scanner := bufio.NewScanner(env.InputStream)
		for scanner.Scan() {
			line++
			fmt.Fprintf(env.OutputStream, "%d.\t%s\n", line, scanner.Text())
		}
	} else {
		io.Copy(env.OutputStream, env.InputStream)
	}
	//	files specified as args
	for _, f := range k.files {
		defer f.Close()
		if k.withNumbering {
			//	output with numbers
			//	ex: 1. hello world
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line++
				fmt.Fprintf(env.OutputStream, "%d.\t%s\n", line, scanner.Text())
			}
		} else {
			_, err := io.Copy(env.OutputStream, f)
			if err != nil {
				lastKnownError = err
				//	should we panic and die, or continue on?
				//	let's be nice and continue on.
				fmt.Fprintln(env.ErrorStream, err)
			}
		}
	}
	return lastKnownError
}
