package platoon_test

import (
	"flag"
	"maps"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/sean9999/go-platoon"
)

type PlatoonError struct {
	msg string
}

func (pe PlatoonError) Error() string {
	return pe.msg
}

func NewPlatoonError(msg string) PlatoonError {
	return PlatoonError{msg}
}

var ErrBingNotAllowed = NewPlatoonError("bing not allowed")

func TestArgParser(t *testing.T) {

	//	parses flags foo and bar, and barfs if there is a flag called bing
	fooAndBarButNotBing := func(args []string) (map[string]any, []string, error) {
		m := map[string]any{}
		var err error
		ts := flag.NewFlagSet("test", flag.ContinueOnError)
		ts.Func("foo", "foo is foo", func(s string) error {
			m["foo"] = s
			return nil
		})
		ts.Func("bar", "bar is bar", func(s string) error {
			m["bar"] = s
			return nil
		})
		err = ts.Parse(args)
		if slices.Contains(args, "bing") {
			err = ErrBingNotAllowed
		}
		return m, ts.Args(), err
	}

	type row struct {
		args        []string
		fn          platoon.ArgumentParser
		expectMarg  map[string]any
		expectTail  []string
		expectError error
	}

	expectMarg1 := map[string]any{
		"foo": "bing",
		"bar": "bat",
	}

	table := []row{
		{[]string{"--foo=bing", "--bar=bat"}, fooAndBarButNotBing, expectMarg1, nil, nil},
		{[]string{"--foo=bing", "--bar=bat", "barbie"}, fooAndBarButNotBing, expectMarg1, []string{"barbie"}, nil},
		{[]string{"--foo=bing", "--bar=bat", "--jazz=pumkin"}, fooAndBarButNotBing, expectMarg1, nil, PlatoonError{}},
	}

	for _, row := range table {
		gotMarg, gottail, gotErr := row.fn(row.args)
		if !maps.Equal(gotMarg, expectMarg1) {
			t.Error(gotMarg)
		}
		if !slices.Equal(row.expectTail, gottail) {
			t.Errorf("expected %v but got %v", row.expectTail, gottail)
		}
		if gotErr != nil {

			if row.expectError == nil {
				t.Errorf("expected no error but got %q", gotErr)
			}

		} else if row.expectError != nil {
			t.Errorf("expected error of type %T but error was nil", row.expectError)
		}
	}

}

func TestArgParser_echo(t *testing.T) {

	//	parses flags as key=value pairs permissively until a non key=value token is found
	allFlagsUntilArgs := func(args []string) (margs, []string, error) {
		m := map[string]any{}
		var err error
		j := 0
		for i, slug := range args {
			kv := strings.Split(slug, "=")
			if len(kv) > 1 {
				reg := regexp.MustCompile(`^\-*`)
				cleanKey := reg.ReplaceAllString(kv[0], "")
				m[cleanKey] = kv[1]
			} else {
				break
			}
			j = i + 1
		}
		return m, args[j:], err
	}

	type row struct {
		args        []string
		fn          platoon.ArgumentParser
		expectMarg  map[string]any
		expectTail  []string
		expectError error
	}

	emptyMap := map[string]any{}

	expectMarg1 := map[string]any{
		"foo": "bing",
		"bar": "bat",
	}

	expectMargWithJazz := map[string]any{
		"foo":  "bing",
		"bar":  "bat",
		"jazz": "pumkin",
	}

	table := []row{
		{[]string{"--foo=bing", "--bar=bat"}, allFlagsUntilArgs, expectMarg1, nil, nil},
		{[]string{"--foo=bing", "--bar=bat", "barbie"}, allFlagsUntilArgs, expectMarg1, []string{"barbie"}, nil},
		{[]string{"barbie"}, allFlagsUntilArgs, emptyMap, []string{"barbie"}, nil},
		{[]string{}, allFlagsUntilArgs, emptyMap, nil, nil},
		{[]string{"--foo=bing", "--bar=bat", "--jazz=pumkin"}, allFlagsUntilArgs, expectMargWithJazz, nil, nil},
	}

	for _, row := range table {
		gotMarg, gottail, gotErr := row.fn(row.args)
		if !maps.Equal(gotMarg, row.expectMarg) {
			t.Errorf("expected margs to be %v, but got %v", row.expectMarg, gotMarg)
		}
		if !slices.Equal(row.expectTail, gottail) {
			t.Errorf("expected tail to be %v, but got %v", row.expectTail, gottail)
		}
		if gotErr != nil {

			if row.expectError == nil {
				t.Error("3", gotErr)
			}

		} else if row.expectError != nil {
			t.Errorf("expected error of type %T but got %T", row.expectError, gotErr)
		}
	}

}
