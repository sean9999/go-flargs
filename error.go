package platoon

import "fmt"

// PlatoonError should be used for all errors in this package
// It can (and should) wrap other errors
type PlatoonError struct {
	msg   string
	child error
}

func (pe *PlatoonError) Error() string {
	if pe.child == nil {
		return pe.msg
	} else {
		return fmt.Sprintf("%s: %q", pe.msg, pe.child)
	}
}

func (pe *PlatoonError) Unwrap() error {
	return pe.child
}

func NewPlatoonError(msg string, childError error) *PlatoonError {
	pe := PlatoonError{
		msg:   msg,
		child: childError,
	}
	return &pe
}
