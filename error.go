package flargs

import "fmt"

// FlargError should be used for all errors in this package
// It can (and should) wrap other errors
type FlargError struct {
	msg   string
	child error
}

func (pe *FlargError) Error() string {
	if pe.child == nil {
		return pe.msg
	} else {
		return fmt.Sprintf("%s: %q", pe.msg, pe.child)
	}
}

func (pe *FlargError) Unwrap() error {
	return pe.child
}

func NewFlargError(msg string, childError error) *FlargError {
	pe := FlargError{
		msg:   msg,
		child: childError,
	}
	return &pe
}
