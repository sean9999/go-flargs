package flargs

// Phase represents lifecycle phases
type Phase uint8

const (
	Uninitialized Phase = iota
	Parsing
	Loading
	Running
	Ran
)
