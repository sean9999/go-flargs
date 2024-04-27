package main

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/sean9999/go-flargs"
)

func TestKat_NormalMode_TwoFiles(t *testing.T) {

	//	fake filesystem
	helloTxt := fstest.MapFile{
		Data: []byte("hello\n"),
	}
	worldTxt := fstest.MapFile{
		Data: []byte("world\n"),
	}
	tmpFs := fstest.MapFS{
		"hello.txt": &helloTxt,
		"world.txt": &worldTxt,
	}

	//	args and flags
	flarguments := []string{
		"hello.txt",
		"world.txt",
	}

	//	parse flargs
	params := new(KatConf)
	_, err := params.Parse(flarguments)
	if err != nil {
		panic(err)
	}

	//	no need for randomness
	env := flargs.NewTestingEnvironment(nil)

	//	add fake fs
	env.Filesystem = tmpFs

	//	try to turn arguments into real file pointers
	err = params.Load(env)
	if err != nil {
		panic(err)
	}

	//	kat
	err = flargs.NewCommand(env, KatFunction).Run(params)

	if err != nil {
		t.Fatal(err)
	}

	//	what is the expected output?
	want := "hello\nworld\n"

	//	what is the actual output?
	got := string(env.GetOutput())

	if want != got {
		t.Error(got)
	}

}

func TestKat_NormalMode_PipedIn(t *testing.T) {

	flarguments := []string{}

	params := new(KatConf)
	_, err := params.Parse(flarguments)
	if err != nil {
		t.Fatalf("error in parsing stage: %s", err)
	}

	env := flargs.NewTestingEnvironment(nil)
	err = params.Load(env)
	if err != nil {
		t.Errorf("error in loading stage: %s", err)
		t.FailNow()
	}

	want := []byte("all your base are belong to us\n")
	env.InputStream.Write(want)

	//	kat
	err = flargs.NewCommand(env, KatFunction).Run(params)
	if err != nil {
		t.Errorf("error in execution stage: %s", err)
		t.FailNow()
	}
	got := env.GetOutput()

	if bytes.Equal(got, want) {
		t.Errorf("wanted %s but got %s", want, got)
	}

}

func TestKat_NumberedMode_TwoFiles(t *testing.T) {
	//	fake filesystem
	helloTxt := fstest.MapFile{
		Data: []byte("hello\n"),
	}
	worldTxt := fstest.MapFile{
		Data: []byte("world\n"),
	}
	tmpFs := fstest.MapFS{
		"hello.txt": &helloTxt,
		"world.txt": &worldTxt,
	}

	//	args and flags
	flarguments := []string{
		"-n",
		"hello.txt",
		"world.txt",
	}

	//	parse flargs
	params := new(KatConf)
	_, err := params.Parse(flarguments)
	if err != nil {
		panic(err)
	}

	//	no need for randomness
	env := flargs.NewTestingEnvironment(nil)

	//	add fake fs
	env.Filesystem = tmpFs

	//	try to turn arguments into real file pointers
	err = params.Load(env)
	if err != nil {
		panic(err)
	}

	//	kat
	err = flargs.NewCommand(env, KatFunction).Run(params)

	if err != nil {
		t.Fatal(err)
	}

	//	what is the expected output?
	want := "1.\thello\n2.\tworld\n"

	//	what is the actual output?
	got := string(env.GetOutput())

	if want != got {
		t.Error(got)
	}

}
