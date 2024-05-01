package flargs

// a Flarger is a custom object that represents the state and functionality of your command
type Flarger interface {
	Parse([]string) error
	Load(*Environment) error
	Run(*Environment) error
}

// StateMachine implements [Flarger]
// it provides basic functionality and default methods
// allowing you to not bother writing them if you don't need them
type StateMachine struct {
	RemainingArgs []string
	Phase         Phase
}

// no-op. This will run if you don't define Parse() in your konf.
func (s *StateMachine) Parse(a []string) error {
	s.Phase = Parsing
	s.RemainingArgs = a
	return nil
}

// no-op. This will run if you don't define Load() in your konf.
func (s *StateMachine) Load(_ *Environment) error {
	s.Phase = Loading
	return nil
}

// no-op. This will run if you don't define Run() in your konf.
func (s *StateMachine) Run(_ *Environment) error {
	s.Phase = Running
	return nil
}
