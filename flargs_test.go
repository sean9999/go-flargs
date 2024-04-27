package flargs_test

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/sean9999/go-flargs"
)

var errBadFileName error = errors.New("bad file name")

type conf struct {
	colour   string
	number   float64
	fileName string
	file     fs.File
}

func (c *conf) Parse(args []string) ([]string, error) {

	fset := flag.NewFlagSet("flags", flag.ContinueOnError)
	fset.Func("colour", "favourite colour", func(s string) error {
		//	should be all lowercase and one of red,green,yellow,etc...
		lstr := strings.ToLower(s)
		switch lstr {
		case "red", "green", "blue", "brown", "yellow", "purple", "orange":
			c.colour = lstr
		default:
			return errors.New("unsupported colour")
		}
		return nil
	})

	fset.Func("number", "floating point number", func(s string) error {
		//	returns error if string is not convertable to float
		n, err := strconv.ParseFloat(s, 64)
		if err == nil {
			c.number = n
		}
		return err
	})
	fset.Parse(args)
	remainders := fset.Args()

	if len(remainders) < 1 {
		return remainders, errors.New("a filename argument is needed")
	}

	c.fileName = remainders[0]

	return remainders[1:], nil
}
func (c *conf) Load(env flargs.Environment) error {
	fd, err := env.Filesystem.Open(c.fileName)
	if err != nil {
		return fmt.Errorf("file %q could not be opened (%w)", c.fileName, errBadFileName)
	}
	c.file = fd
	return nil
}

func TestFlargs(t *testing.T) {

	goMod, _ := os.Open("go.mod")

	type row struct {
		inputArgs  []string
		conf       conf
		remainders []string
		err        error
	}
	table := []row{
		{ //	test #1
			inputArgs:  []string{"--colour=red", "--number=7", "go.mod"},
			conf:       conf{"red", 7, "go.mod", goMod},
			remainders: nil,
			err:        nil,
		},
		{ //	test #2
			inputArgs:  []string{"--colour=red", "--number=7", "go.mod", "fish", "BBQ"},
			conf:       conf{"red", 7, "go.mod", goMod},
			remainders: []string{"fish", "BBQ"},
			err:        nil,
		},
		{
			inputArgs:  []string{"--colour=blue", "--number=7.17", "go.mod", "fish", "BBQ"},
			conf:       conf{"blue", 7.17, "go.mod", goMod},
			remainders: []string{"fish", "BBQ"},
			err:        nil,
		},
		{
			inputArgs:  []string{"--colour=blue", "--number=7.17", "go.mod_x", "fish", "BBQ"},
			conf:       conf{"blue", 7.17, "go.mod", goMod},
			remainders: []string{"fish", "BBQ"},
			err:        errBadFileName,
		},
	}

	for _, want := range table {

		got := new(conf)
		gotRemainders, gotErr := got.Parse(want.inputArgs)

		if gotErr != nil && want.err != nil {
			if !errors.Is(gotErr, want.err) {
				t.Fatal(gotErr)
			}
		}

		if !slices.Equal(want.remainders, gotRemainders) {
			//t.Error(gotRemainders)
			t.Errorf("expected %v but got %v", want.remainders, gotRemainders)
		}

		if got.colour != want.conf.colour {
			t.Errorf("wanted %q but got %q", want.conf.colour, got.colour)
		}
		if got.number != want.conf.number {
			t.Errorf("wanted %f but got %f", want.conf.number, got.number)
		}

		wantFileName := want.conf.fileName
		gotFileName := goMod.Name()
		if gotFileName != wantFileName {
			t.Errorf("wanted %q but got %q", wantFileName, gotFileName)
		}

	}

}

// an object that represents the input you need
type helloConf struct {
	name string
}

func (c *helloConf) Parse(args []string) (remainders []string, err error) {
	//  default value
	c.name = "world"
	fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
	fset.Func("name", "hello to who?", func(s string) error {
		if s == "batman" {
			return errors.New("you cannot say hello to batman")
		}
		c.name = s
		return nil
	})
	err = fset.Parse(args)
	return fset.Args(), err
}
func (c *helloConf) Load(env *flargs.Environment) error {
	return nil
}

func ExampleNewCommand_hello() {

	// this is a flargs.RunFunc. It says hello.
	helloFn := func(env *flargs.Environment, conf *helloConf) error {
		outputString := fmt.Sprintf("hello, %s", conf.name)
		env.OutputStream.Write([]byte(outputString))
		return nil
	}

	conf := new(helloConf)

	conf.Parse([]string{"--name", "robin"})

	env := flargs.NewTestingEnvironment(nil)
	cmd := flargs.NewCommand(env, helloFn)
	cmd.Run(conf)

	got := string(cmd.Env.GetOutput())

	fmt.Println(got)
	// Output: hello, robin

}
