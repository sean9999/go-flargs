package main

import (
	"math/rand"
	"testing"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
)

func TestMain(t *testing.T) {

	t.Run("with no params", func(t *testing.T) {

		params := new(proverbs.Params)
		env := flargs.NewTestingEnvironment(rand.NewSource(0))
		cmd := flargs.NewCommand(params, env)
		err := cmd.ParseAndLoad([]string{})
		if err != nil {
			panic(err)
		}
		cmd.Run()
		want := "Don't panic.\n"
		got := string(env.GetOutput())
		if want != got {
			t.Errorf("wanted %s but got %s", want, got)
		}

	})

}
