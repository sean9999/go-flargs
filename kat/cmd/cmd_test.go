package main

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/kat"
)

func TestKat_with_one_arg(t *testing.T) {

	want := []byte("all your base are belong to us.")

	konf := new(kat.Konf)
	env := flargs.NewTestingEnvironment(nil)
	tmpFs := fstest.MapFS{
		"base.txt": &fstest.MapFile{
			Data: want,
		},
	}
	env.Filesystem = tmpFs

	katCmd := flargs.NewCommand(konf, env)
	inputParams := []string{
		"base.txt",
	}

	err := katCmd.ParseAndLoad(inputParams)
	if err != nil {
		t.Fatal(err)
	}
	err = katCmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	got := env.GetOutput()
	if !bytes.Equal(got, want) {
		t.Errorf("got %s but wanted %s", got, want)
	}

}

func TestKat_stdin(t *testing.T) {

	state := new(kat.Konf)
	want := []byte("all your base are belong to us.")
	env := flargs.NewTestingEnvironment(nil)

	katCmd := flargs.NewCommand(state, env)

	// no arguments. We're passing in data to stdin
	env.InputStream.Write(want)
	err := katCmd.ParseAndLoad(nil)

	if err != nil {
		t.Error(err)
	}

	katCmd.Run()

	got := env.GetOutput()
	if !bytes.Equal(got, want) {
		t.Errorf("got %s but wanted %s", got, want)
	}

}
