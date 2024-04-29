package flargs_test

import (
	"bytes"
	"testing"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/kat"
)

func TestKatOneArg(t *testing.T) {

	noInput := []string{}
	want := []byte("all your base are belong to us")

	env := flargs.NewTestingEnvironment(nil)

	konf := new(kat.Konf)

	cmd := flargs.NewCommand(konf, env)
	env.InputStream.Write(want)

	err := cmd.ParseAndLoad(noInput)
	if err != nil {
		t.Error(err)
	}

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	got := env.GetOutput()

	if !bytes.Equal(want, got) {
		t.Errorf("got %s but wanted %s", got, want)
	}

}
