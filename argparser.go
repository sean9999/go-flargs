package platoon

// margs is a map of args
type margs = map[string]any

// ArgumentParser takes in a slice of strings and returns something useful to a CommandFunction
type ArgumentParser func([]string) (margs, []string, error)

// type parseMachine struct {
// 	cursor int
// 	args   []string
// 	margs  map[string]any
// 	tail   []string
// }
