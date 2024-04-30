package flargs_test

import (
	"bytes"
	"testing"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/kat"
)

func TestKat(t *testing.T) {

	want := []byte("all your base are belong to us")
	env := flargs.NewTestingEnvironment(nil)
	konf := new(kat.Konf)
	cmd := flargs.NewCommand(konf, env)

	//	pipe to stdin
	env.InputStream.Write(want)

	//	no flags. no args
	err := cmd.ParseAndLoad([]string{})
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
