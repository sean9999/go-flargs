package kat

import (
	"flag"
	"io/fs"
)

// Konf is all the information necessary to parse flags for, and then prepare to execute kat
type Konf struct {
	fileNames     []string  // filled at parse-time
	files         []fs.File // created at command-time
	withNumbering bool
	remainingArgs []string
}

func (k *Konf) Args() []string {
	return k.remainingArgs
}

// do these flags make sense?
// an error will result if any flag other than "-n" is passed in
// all arguments are assumed to be files.
// We don't have access to an Environment here, so this logic
// represents all files that are *theoretically* possible
func (k *Konf) Parse(args []string) error {
	fset := flag.NewFlagSet("flargs", flag.ContinueOnError)
	fset.BoolVar(&k.withNumbering, "n", false, "use numbering")
	err := fset.Parse(args)
	if err != nil {
		k.remainingArgs = fset.Args()
		return err
	}
	k.fileNames = fset.Args()
	//	there should be no remaining args because we're consuming them all
	k.remainingArgs = []string{}
	return nil
}
