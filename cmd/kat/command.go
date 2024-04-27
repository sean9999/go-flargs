package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/sean9999/go-flargs"
)

// the kat function
var KatFunction = func(env *flargs.Environment, input *KatConf) error {
	var lastKnownError error
	line := 0

	//	stdin
	if input.withNumbering {
		scanner := bufio.NewScanner(env.InputStream)
		for scanner.Scan() {
			line++
			fmt.Fprintf(env.OutputStream, "%d.\t%s\n", line, scanner.Text())
		}
	} else {
		io.Copy(env.InputStream, env.OutputStream)
	}

	//	files specified as args
	for _, f := range input.files {
		defer f.Close()
		if input.withNumbering {

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
