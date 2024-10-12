package executing

// Executor is an interface that defines basic operations for executing a command.
type Executor interface {
	Execute(options map[string]bool) error
}