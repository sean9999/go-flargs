package flargs

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
