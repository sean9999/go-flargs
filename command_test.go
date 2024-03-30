package platoon_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"slices"
	"testing"

	"github.com/sean9999/go-platoon"
)

type margs = map[string]any

func TestNewCommand(t *testing.T) {

	env := &platoon.Environment{
		InputStream:  new(bytes.Buffer),
		OutputStream: new(bytes.Buffer),
		ErrorStream:  new(bytes.Buffer),
		Variables: map[string]string{
			"RUN_CONTEXT": "go test",
		},
	}

	//	the command
	fn := func(env *platoon.Environment, margs map[string]any) error {
		j, err := json.Marshal(margs)
		if err != nil {
			return err
		}
		env.OutputStream.Write(j)
		return nil
	}

	cmd := platoon.NewCommand(env, fn)

	//	the margs
	m := map[string]any{
		"foo": "bar",
		"bat": "bing",
	}

	want, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}

	//	run the command with margs
	err = cmd.Run(m)
	if err != nil {
		t.Error(err)
	}

	got := bytes.NewBuffer(nil)
	got.ReadFrom(env.OutputStream)

	if !slices.Equal(want, got.Bytes()) {
		t.Error(got.String(), "\t", string(want))
	}

}

func Example_subcommand() {

	//	let's implement a utility with a subcommand like this
	//	busybox git --work-tree=/some/work/tree checkout --branch=main https://github.com/sean9999/go-platoon

	//	if first subcommand is not busybox, die
	//	if second subcommand is not "git", exit with "unsupported" warning

	//	git takes a margs that may include a "work-tree" prop, or it may simply be empty
	//	it should return ["checkout", ...] as a tail

	//	checkout sub-subcommand should have a margs of "branch" and "repo"
	//	repo is not a flag. It's an arg.

	type margs = map[string]any
	inputParams := []string{"busybox", "git", "--work-tree=/some/folder", "checkout", "https://github.com/sean9999/go-platoon"}

	rootParser := func(tokens []string) (margs, []string, error) {
		if tokens[0] != "busybox" {
			err := platoon.NewPlatoonError("first command is not busybox", nil)
			return nil, nil, err
		}
		tail := tokens[1:]
		return nil, tail, nil
	}

	busyBoxMarg, busyboxRemainder, err := rootParser(inputParams)

	fmt.Println("busybox", busyBoxMarg)
	fmt.Println("busybox", busyboxRemainder)
	fmt.Println("busybox", err)

	gitParser := func(tokens []string) (margs, []string, error) {
		m := map[string]any{}
		if tokens[0] != "git" {
			return nil, nil, platoon.NewPlatoonError(fmt.Sprintf("subcommand %q is not %q", tokens[0], "git"), nil)
		}
		fset := *flag.NewFlagSet("git", flag.ContinueOnError)
		fset.Func("work-tree", "git work tree", func(s string) error {
			m["workTree"] = s
			return nil
		})
		fset.Parse(tokens[1:])
		return m, fset.Args(), nil
	}

	gitMarg, gitsubmoduleParams, err := gitParser(busyboxRemainder)

	fmt.Println("git", gitMarg)
	fmt.Println("git", gitsubmoduleParams)
	fmt.Println("git", err)

	// Output:
	// busybox map[]
	// busybox [git --work-tree=/some/folder checkout https://github.com/sean9999/go-platoon]
	// busybox <nil>
	// git map[workTree:/some/folder]
	// git [checkout https://github.com/sean9999/go-platoon]
	// git <nil>

}
