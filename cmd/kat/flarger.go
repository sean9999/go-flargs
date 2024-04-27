package main

import (
	"flag"
	"io/fs"

	"github.com/sean9999/go-flargs"
)

type KatConf struct {
	fileNames     []string  // filled at parse-time
	files         []fs.File // created at command-time
	withNumbering bool
}

func (k *KatConf) Parse(args []string) ([]string, error) {
	fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
	fset.BoolVar(&k.withNumbering, "n", false, "use numbering")
	err := fset.Parse(args)
	if err != nil {
		return fset.Args(), err
	}
	k.fileNames = fset.Args()
	//	there should be no remaining args because we're consuming them all
	return []string{}, nil
}
func (k *KatConf) Load(env *flargs.Environment) error {
	k.files = []fs.File{}
	for _, fileName := range k.fileNames {
		fd, err := env.Filesystem.Open(fileName)
		if err != nil {
			return err
		}
		k.files = append(k.files, fd)
	}
	return nil
}
