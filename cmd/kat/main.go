package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
)

func main() {

	//	input object suitable for passing into cat
	type katKonf struct {
		files         []*os.File
		withNumbering bool
	}

	//	argument parser for katFn
	var katFlarger flargs.ParseFunc[*katKonf] = func(args []string) (*katKonf, []string, error) {
		conf := new(katKonf)
		conf.files = []*os.File{}
		fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
		fset.BoolVar(&conf.withNumbering, "n", false, "use numbering")
		err := fset.Parse(args)
		if err != nil {
			return nil, fset.Args(), err
		}
		if len(fset.Args()) < 1 {
			fi, _ := os.Stdin.Stat()
			//	if data was piped in, we don't need to panic on zero args
			if fi.Size() == 0 {
				return nil, fset.Args(), errors.New("must specify at least one file")
			} else {
				conf.files = append(conf.files, os.Stdin)
			}
		}
		for _, arg := range fset.Args() {
			f, e := os.Open(arg)
			if e != nil {
				return nil, fset.Args(), fmt.Errorf("could not open %q: %w", arg, e)
			}
			conf.files = append(conf.files, f)
		}
		//	there should be no remaining args because we're consuming them all
		return conf, []string{}, nil
	}

	//	the kat function
	katFn := func(env *flargs.Environment, input *katKonf) error {
		var lastKnownError error
		line := 0
		for _, f := range input.files {
			defer f.Close()

			if input.withNumbering {
				scanner := bufio.NewScanner(f)

				for scanner.Scan() {
					line++
					fmt.Fprintf(env.OutputStream, "%d.\t%s\n", line, scanner.Text())
				}
			} else {
				_, err := f.WriteTo(env.OutputStream)
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

	env := flargs.NewCLIEnvironment()
	cmd := flargs.NewCommand(env, katFn)
	params, _, err := katFlarger(os.Args[1:])

	if err != nil {
		panic(err)
	}

	cmd.Run(params)

}
