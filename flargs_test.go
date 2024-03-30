package platoon_test

import (
	"errors"
	"flag"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/sean9999/go-platoon"
)

type row struct {
	inputs          []string
	parseFunc       platoon.FlargsParseFunc
	expectMap       map[string]any
	expectRemainder []string
	expectError     error
}

var ErrInvalidFlag error = platoon.NewPlatoonError("platoon error", errors.New("invalid flag"))

func TestNewFlargs(t *testing.T) {

	var ErrPlatoon *platoon.PlatoonError

	parseFn := func(args []string) (map[string]any, []string, error) {
		margs := map[string]any{
			"bing": "PINEAPPLE", // set default values here
		}

		fset := new(flag.FlagSet)
		fset.Func("foo", "foo dee doo", func(s string) error {
			margs["foo"] = s
			return nil
		})
		fset.Func("bar", "a bar is a a bar", func(s string) error {
			//	must begin with the letter "b"
			ss := strings.Split(s, "")
			if ss[0] != "b" && ss[0] != "B" {
				return ErrInvalidFlag
			}
			margs["bar"] = s
			return nil
		})
		fset.Func("bing", "bing is usually a pineapple", func(s string) error {
			margs["bing"] = s
			return nil
		})

		//	this will hydrate margs
		err := fset.Parse(args)

		if err != nil {
			err = fmt.Errorf("arg parse error: %w", err)
			//err = platoon.NewPlatoonError("arg parse error", err)
		}

		//	further validation or massaging as needed

		return margs, fset.Args(), err
	}

	table := []row{
		{
			inputs:    []string{"--foo=bonk", "--bar=bink"},
			parseFunc: parseFn,
			expectMap: map[string]any{
				"foo":  "bonk",
				"bar":  "bink",
				"bing": "PINEAPPLE",
			},
			expectRemainder: []string{},
			expectError:     nil,
		},
		{
			inputs:    []string{"--foo=bonk", "--bar=TALYHOE"},
			parseFunc: parseFn,
			expectMap: map[string]any{
				"foo":  "bonk",
				"bing": "PINEAPPLE",
			},
			expectRemainder: []string{},
			expectError:     ErrPlatoon,
		},
		{
			inputs:    []string{"--foo=bonk", "--bar=bink", "--bing=bank"},
			parseFunc: parseFn,
			expectMap: map[string]any{
				"foo":  "bonk",
				"bar":  "bink",
				"bing": "bank",
			},
			expectRemainder: []string{},
			expectError:     nil,
		},
		{
			inputs:    []string{"--foo=bonk", "--bar=bink", "darth", "vader"},
			parseFunc: parseFn,
			expectMap: map[string]any{
				"foo":  "bonk",
				"bar":  "bink",
				"bing": "PINEAPPLE",
			},
			expectRemainder: []string{"darth", "vader"},
			expectError:     nil,
		},
	}

	for _, row := range table {

		//fset := flag.NewFlagSet("fset", flag.ContinueOnError)
		fl := platoon.NewFlargs(parseFn)
		m, remainder, err := fl.Parse(row.inputs)

		if err != nil {

			if row.expectError != nil {

				if !errors.As(err, &ErrPlatoon) {
					t.Error("erros.Is thing was false ", err)
				}

			}

		}

		if !slices.Equal(remainder, row.expectRemainder) {
			t.Errorf("remainder was %v", remainder)
		}

		if !maps.Equal(m, row.expectMap) {
			t.Errorf("the map we got was %v, but we wanted %v", m, row.expectMap)
		}

	}

}

func Example_foobar() {

	//	parse foo and bar. Return the remaining args as a tail
	fooBargs := platoon.NewFlargs(func(args []string) (map[string]any, []string, error) {
		m := map[string]any{}
		fs := new(flag.FlagSet)
		fs.Func("foo", "foo must be a non-empty string", func(s string) error {
			if len(s) == 0 {
				return errors.New("foo cannot be zero-length")
			}
			m["foo"] = s
			return nil
		})
		fs.Func("bar", "bar must be a number", func(s string) error {
			i, err := strconv.Atoi(s)
			if err == nil {
				m["bar"] = i
			}
			return err
		})
		fs.Parse(args)
		return m, fs.Args(), nil
	})

	//	call the above function with command-line args
	margs, tail, err := fooBargs.Parse([]string{"--foo=HELLO", "--bar=79", "mysubcommand", "myarg"})

	fmt.Println(margs["foo"])
	fmt.Println(margs["bar"])
	fmt.Println(tail)
	fmt.Println(err)

	// Output:
	// HELLO
	// 79
	// [mysubcommand myarg]
	// <nil>

}
