package main

import (
	"bytes"
	"testing"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/rot13"
)

func TestRot13(t *testing.T) {

	inputText := []byte("neon penny nowhere germ, pening roof balk.")
	want := []byte("arba craal abjurer trez, cravat ebbs onyx.")

	k := new(rot13.RotKonf)
	env := flargs.NewTestingEnvironment(nil) //  no randomness needed
	cmd := flargs.NewCommand(k, env)
	cmd.ParseAndLoad(nil)

	//	send some text to "stdin"
	env.InputStream.Write(inputText)

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

	got := env.GetOutput()

	if !bytes.Equal(got, want) {
		t.Errorf("wanted %s but got %s", want, got)
	}

}
