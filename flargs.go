package flargs

type margs = map[string]any

// Flargs parses a slice of strings into a well-formed map and returns the remaining unparsed arguments.
// It exists solely to capture errors from it's [FlargsParseFunc] and wrap them in a [FlargError].
type Flargs interface {
	Parse([]string) (map[string]any, []string, error)
}

type flargs struct {
	ParseFunc FlargsParseFunc
}

type FlargsParseFunc func([]string) (map[string]any, []string, error)

func (f *flargs) Parse(args []string) (map[string]any, []string, error) {
	m, remainder, err := f.ParseFunc(args)
	//	wrap any error in a PlatoonError for easy identification
	if err != nil {
		err = &FlargError{"error parsing flags", err}
	}
	return m, remainder, err
}

func NewFlargs(fn FlargsParseFunc) Flargs {
	fl := flargs{fn}
	return &fl
}
