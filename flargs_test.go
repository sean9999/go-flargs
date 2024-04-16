package flargs_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/sean9999/go-flargs"
)

var errBadFileName error = errors.New("bad file name")

type conf struct {
	colour string
	number float64
	file   *os.File
}

func TestFlargs(t *testing.T) {

	goMod, _ := os.Open("go.mod")

	var fparse flargs.ParseFunc[*conf] = func(args []string) (*conf, []string, error) {
		conf := new(conf)
		fset := flag.NewFlagSet("flags", flag.ContinueOnError)
		fset.Func("colour", "favourite colour", func(s string) error {
			//	should be all lowercase and one of red,green,yellow,etc...
			lstr := strings.ToLower(s)
			switch lstr {
			case "red", "green", "blue", "brown", "yellow", "purple", "orange":
				conf.colour = lstr
			default:
				return errors.New("unsupported colour")
			}
			return nil
		})

		fset.Func("number", "floating point number", func(s string) error {
			//	returns error if string is not convertable to float
			n, err := strconv.ParseFloat(s, 64)
			if err == nil {
				conf.number = n
			}
			return err
		})
		fset.Parse(args)
		remainders := fset.Args()

		if len(remainders) < 1 {
			return nil, remainders, errors.New("a filename argument is needed")
		}

		filePath := remainders[0]
		fd, err := os.Open(filePath)
		if err != nil {
			return nil, remainders[1:], fmt.Errorf("file %q could not be opened (%w)", filePath, errBadFileName)
		}
		conf.file = fd

		return conf, remainders[1:], err
	}

	type row struct {
		inputArgs  []string
		conf       conf
		remainders []string
		err        error
	}
	table := []row{
		{ //	test #1
			inputArgs:  []string{"--colour=red", "--number=7", "go.mod"},
			conf:       conf{"red", 7, goMod},
			remainders: nil,
			err:        nil,
		},
		{ //	test #2
			inputArgs:  []string{"--colour=red", "--number=7", "go.mod", "fish", "BBQ"},
			conf:       conf{"red", 7, goMod},
			remainders: []string{"fish", "BBQ"},
			err:        nil,
		},
		{
			inputArgs:  []string{"--colour=blue", "--number=7.17", "go.mod", "fish", "BBQ"},
			conf:       conf{"blue", 7.17, goMod},
			remainders: []string{"fish", "BBQ"},
			err:        nil,
		},
		{
			inputArgs:  []string{"--colour=blue", "--number=7.17", "go.mod_x", "fish", "BBQ"},
			conf:       conf{"blue", 7.17, goMod},
			remainders: []string{"fish", "BBQ"},
			err:        errBadFileName,
		},
	}

	for _, want := range table {
		got, gotRemainders, gotErr := fparse(want.inputArgs)

		if gotErr != nil && want.err != nil {
			if !errors.Is(gotErr, want.err) {
				t.Fatal(gotErr)
			}
		}

		if !slices.Equal(want.remainders, gotRemainders) {
			//t.Error(gotRemainders)
			t.Errorf("expected %v but got %v", want.remainders, gotRemainders)
		}

		if got != nil {
			if got.colour != want.conf.colour {
				t.Errorf("wanted %q but got %q", want.conf.colour, got.colour)
			}
			if got.number != want.conf.number {
				t.Errorf("wanted %f but got %f", want.conf.number, got.number)
			}
		}

		//	@todo: do inode numbers make more sense here?
		wantFileName := want.conf.file.Name()
		gotFileName := goMod.Name()
		if gotFileName != wantFileName {
			t.Errorf("wanted %q but got %q", wantFileName, gotFileName)
		}
		// if want.conf.file.Fd() != got.file.Fd() {
		// 	t.Errorf("got %d but wanted %d", got.file.Fd(), want.conf.file.Fd())
		// }

	}

}

func ExampleNewCommand_hello() {

	type helloConf struct {
		name string
	}

	// this is a flargs.ParseFunc
	parseFn := func(args []string) (*helloConf, []string, error) {
		conf := new(helloConf)
		//  default value
		conf.name = "world"
		fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
		fset.Func("name", "hello to who?", func(s string) error {
			if s == "batman" {
				return errors.New("you cannot say hello to batman")
			}
			conf.name = s
			return nil
		})
		err := fset.Parse(args)
		return conf, fset.Args(), err
	}

	// this is a flargs.RunFunc
	helloFn := func(env *flargs.Environment, conf *helloConf) error {
		outputString := fmt.Sprintf("hello, %s", conf.name)
		env.OutputStream.Write([]byte(outputString))
		return nil
	}

	conf, _, _ := parseFn([]string{"--name", "robin"})

	env := flargs.NewTestingEnvironment(nil)
	cmd := flargs.NewCommand(env, helloFn)
	cmd.Run(conf)

	got := new(bytes.Buffer)
	got.ReadFrom(cmd.Env.OutputStream)

	fmt.Println(got.String())
	// Output: hello, robin

}
