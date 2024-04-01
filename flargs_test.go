package flargs_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/sean9999/go-flargs"
)

type conf struct {
	colour string
	number float64
	file   *os.File
}

func TestFlargs(t *testing.T) {

	goMod, _ := os.Open("go.mod")

	var fparse flargs.FlargsParseFunc[*conf] = func(args []string) (*conf, []string, error) {
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
		nPtr := fset.Float64("number", 0, "favourite numbers")
		// fset.Func("file", "must be a real file", func(s string) error {
		// 	fd, err := os.Open(s)
		// 	if err == nil {
		// 		conf.file = fd
		// 	}
		// 	return err
		// })
		fset.Parse(args)
		conf.number = *nPtr
		remainders := fset.Args()

		if len(remainders) < 1 {
			return nil, nil, errors.New("a filename argument is needed")
		}

		filePath := remainders[0]
		fd, err := os.Open(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("file %q could not be opened", filePath)
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
			err:        nil,
		},
	}

	for _, want := range table {
		got, gotRemainders, gotErr := fparse(want.inputArgs)
		if gotErr != want.err {
			//	A consequence of this fatal error is that following tests will fail.
			//	So skip those in the name of brevity. This error matters most.
			t.Fatal(gotErr)
		}
		if !slices.Equal(want.remainders, gotRemainders) {
			t.Error(gotRemainders)
		}
		if got.colour != want.conf.colour {
			t.Errorf("wanted %q but got %q", want.conf.colour, got.colour)
		}
		if got.number != want.conf.number {
			t.Errorf("wanted %f but got %f", want.conf.number, got.number)
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
