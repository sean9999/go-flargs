package flargs_test

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sean9999/go-flargs"
	"k8s.io/utils/diff"
)

var goWorkUpper = strings.TrimSpace(`
GO 1.22.1

USE .
`)

var humansAndGoWork = strings.TrimSpace(`
/* ME */

Name: Sean Macdonald
Site: https://www.seanmacdonald.ca
go 1.22.1

use .
`)

var humansAndGoWorkNumbered = strings.TrimSpace(`
1.	/* ME */
2.	
3.	Name: Sean Macdonald
4.	Site: https://www.seanmacdonald.ca
5.	go 1.22.1
6.	
7.	use .
`)

func TestNewCommand_cat(t *testing.T) {

	//	input object suitable for passing into cat
	type catConf struct {
		files         []*os.File
		withNumbering bool
	}

	//	argument parser for catFn
	var catFlagParser flargs.ParseFunc[*catConf] = func(args []string) (*catConf, []string, error) {
		conf := new(catConf)
		conf.files = []*os.File{}
		fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
		fset.BoolVar(&conf.withNumbering, "n", false, "use numbering")
		err := fset.Parse(args)
		if err != nil {
			return nil, fset.Args(), err
		}
		if len(fset.Args()) < 1 {
			return nil, fset.Args(), errors.New("must specify at least one file")
		}
		for _, arg := range fset.Args() {
			f, e := os.Open(arg)
			if e != nil {
				return nil, fset.Args(), fmt.Errorf("Could not open %q: %w", arg, e)
			}
			conf.files = append(conf.files, f)
		}
		//	there should be no remaining args because we're consuming them all
		return conf, []string{}, nil
	}

	//	the cat function
	catFn := func(env *flargs.Environment, input *catConf) error {
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

	//	testing *Environment
	env := flargs.NewTestingEnvironment()

	//	catCmd is catFn + env
	catCmd := flargs.NewCommand(env, catFn)

	//	uppercaseify
	//no input needed

	upperFn := func(env *flargs.Environment, _ *struct{}) error {
		plainText := new(bytes.Buffer)
		plainText.ReadFrom(env.InputStream)
		upperText := strings.ToUpper(plainText.String())
		_, err := env.OutputStream.Write([]byte(upperText))
		return err
	}

	upperCmd := flargs.NewCommand(env, upperFn)

	t.Run("cat some files", func(t *testing.T) {

		type col struct {
			inputArgs  []string
			wantErr    error
			wantResult string
		}

		table := []col{
			{[]string{"humans.txt", "go.work"}, nil, humansAndGoWork},
			{[]string{"-n", "humans.txt", "go.work"}, nil, humansAndGoWorkNumbered},
		}

		for _, row := range table {
			konf, _, err := catFlagParser(row.inputArgs)

			if err != nil {
				if !errors.Is(err, row.wantErr) {
					t.Errorf("wanted %s but got %s", row.wantErr, err)
					t.FailNow()
				}
			} else {
				if row.wantErr != nil {
					t.Errorf("wanted error %s but got nil", row.wantErr)
					t.FailNow()
				}
				catCmd.Run(konf)
				//	compare output to expected
				gotBuff := new(bytes.Buffer)
				gotBuff.ReadFrom(catCmd.Env.OutputStream)
				got := strings.TrimSpace(gotBuff.String())
				if row.wantResult != got {
					t.Error(diff.StringDiff(got, row.wantResult))
				}
			}

		}

	})

	t.Run("uppercase-ify", func(t *testing.T) {

		type row struct {
			inputString  string
			expectError  error
			expectString string
		}

		table := []row{
			{"all your base", nil, "ALL YOUR BASE"},
		}

		for _, row := range table {

			upperCmd.Env.InputStream.Write([]byte(row.inputString))
			err := upperCmd.Run(nil)
			if err == nil {

				gotBuf, err := io.ReadAll(upperCmd.Env.OutputStream)
				if err != nil {
					t.Error(err)
				}

				if string(gotBuf) != row.expectString {
					t.Errorf("wanted %q but got %q", row.expectString, gotBuf)
				}

			} else {
				t.Error(err)
			}

		}

	})

	t.Run("pipe cat to uppercase-ify", func(t *testing.T) {

		type row struct {
			inputString  string
			expectError  error
			expectString string
		}

		table := []row{
			{"go.work", nil, goWorkUpper},
		}

		for _, row := range table {

			konf, _, err := catFlagParser([]string{row.inputString})

			if err == nil {
				//	pipe cat to uppcaseify
				catCmd.Pipe(konf, upperCmd.Env)
				upperCmd.Run(nil)

				resultBytes, err := io.ReadAll(upperCmd.Env.OutputStream)
				if err != nil {
					t.Error(err)
				}
				resultString := strings.TrimSpace(string(resultBytes))

				if resultString != row.expectString {
					t.Errorf("was expecting %s but got %s", goWorkUpper, resultString)
				}

			} else {
				t.Error(err)
			}

		}

	})

}
