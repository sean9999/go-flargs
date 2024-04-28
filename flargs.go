package flargs

// a Flarger is a custom object that represents the state and functionality of your command
type Flarger interface {
	Parse([]string) error
	Load(*Environment) error
	Run(*Environment) error
	Args() []string // for remaining args after a parse
}
