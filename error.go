package flargs

import "fmt"

type ExitCode uint8

const (
	ExitCodeSuccess ExitCode = iota
	ExitCodeGenericError
	ExitCodeMisuseOfBuiltIns
	ExitCodeCannotExecute   = iota + 123
	ExitCodeCommandNotFound // Command not found
	ExitCodeInvalidArgumentToExit
	ExitCodeFatalErrorSignal1
	ExitCodeFatalErrorSignal2 // Ctrl-C was pressed
	ExitCodeFatalErrorSignal3
	ExitCodeFatalErrorSignal4
	ExitCodeFatalErrorSignal5
	ExitCodeFatalErrorSignal6
	ExitCodeFatalErrorSignal7
	ExitCodeFatalErrorSignal8
	ExitCodeFatalErrorSignal9
)

func (ec ExitCode) Error() string {
	return fmt.Sprintf("exit code: %d", ec)
}

type FlargError struct {
	ExitCode
	UnderlyingError error
}

func (fe *FlargError) Error() string {
	return fmt.Sprintf("flargs error: %s", fe.UnderlyingError)
}

func NewFlargError(exitcode ExitCode, underlying error) *FlargError {
	fe := &FlargError{exitcode, underlying}
	return fe
}
