package platoon

// a Command is composed of a [CommandFunction] and an execution [Environment].
// You can run it with Run(), which expects arguments in the form map[string]any (margs)
type Command interface {
	Run(margs) error
}

// CommandFunction is the function that executes via Run()
// it should send success data to Environment.OutputStream
// and error data to Environment.ErrorStream
type CommandFunction func(*Environment, margs) error

// command implements Command
type command struct {
	env *Environment
	exe CommandFunction
}

func NewCommand(env *Environment, exe CommandFunction) Command {
	cmd := command{env, exe}
	return &cmd
}

func (com *command) Run(m margs) error {
	return com.exe(com.env, m)
}
