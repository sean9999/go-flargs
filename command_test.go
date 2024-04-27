package flargs_test

import (
	"bufio"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"

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

// input object suitable for passing into kat
type katConf struct {
	fileNames     []string
	files         []fs.File
	withNumbering bool
}

func (c *katConf) Parse(args []string) ([]string, error) {

	fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
	fset.BoolVar(&c.withNumbering, "n", false, "use numbering")
	err := fset.Parse(args)
	if err != nil {
		return fset.Args(), err
	}
	c.fileNames = fset.Args()

	//	there should be no remaining args because kat consumes them all
	return []string{}, nil
}
func (c *katConf) Load(env *flargs.Environment) error {
	c.files = []fs.File{}
	for _, arg := range c.fileNames {
		f, e := env.Filesystem.Open(arg)
		if e != nil {
			return fmt.Errorf("Could not open %q: %w", arg, e)
		}
		c.files = append(c.files, f)
	}
	return nil
}

func TestNewCommand_cat(t *testing.T) {

	//	the cat function
	catFn := func(env *flargs.Environment, input *katConf) error {
		var lastKnownError error
		line := 0

		//	if something was piped in, kat that first
		io.Copy(env.InputStream, env.OutputStream)

		//	now run through files passed in as arguments
		for _, f := range input.files {
			defer f.Close()

			if input.withNumbering {
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

	//	testing Environment
	env := flargs.NewTestingEnvironment(rand.Reader)

	//	create a filesystem with the files we want
	mfs := fstest.MapFS{}
	humansContent, _ := os.ReadFile("humans.txt")
	goWorkContent, _ := os.ReadFile("go.work")
	mfs["humans.txt"] = &fstest.MapFile{
		Data: humansContent,
	}
	mfs["go.work"] = &fstest.MapFile{
		Data: goWorkContent,
	}
	env.Filesystem = mfs

	//	catCmd is catFn + env
	catCmd := flargs.NewCommand(env, catFn)

	//	uppercaseify
	//	no flargs needed. passing nil is ok
	upperFn := func(env *flargs.Environment, _ flargs.Flarger[any]) error {
		plainText := string(env.GetInput())
		upperText := strings.ToUpper(plainText)
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

			konf := new(katConf)

			_, err := konf.Parse(row.inputArgs)

			if err != nil {
				if !errors.Is(err, row.wantErr) {
					t.Errorf("wanted %s but got %s", row.wantErr, err)
					t.FailNow()
				}
			} else {

				//	load only runs if Parse succeeded
				err := konf.Load(env)

				if row.wantErr != nil && err == nil {
					t.Errorf("wanted error %s but got nil", row.wantErr)
					t.FailNow()
				}
				catCmd.Run(konf)

				//	compare output to expected
				got := strings.TrimSpace(string(catCmd.Env.GetOutput()))
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
			konf := new(katConf)
			_, err := konf.Parse([]string{row.inputString})
			if err == nil {

				konf.Load(env)
				//	pipe cat to uppcaseify
				catCmd.Pipe(konf, upperCmd.Env)
				err = catCmd.Run(konf)
				if err != nil {
					t.Error(err)
				}
				upperCmd.Run(nil)
				result := strings.TrimSpace(string(upperCmd.Env.GetOutput()))
				if result != row.expectString {
					t.Errorf("was expecting %s but got %s", goWorkUpper, result)
				}

			} else {
				t.Error(err)
			}
		}
	})

}
