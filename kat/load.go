package kat

import (
	"io/fs"

	"github.com/sean9999/go-flargs"
)

// now that we have access to Environment
// let's see if all the strings passed in can be mapped to real files
// return an error for any non-valid file name.
func (k *Konf) Load(env *flargs.Environment) error {
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
