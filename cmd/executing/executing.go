package executing

type Executor interface {
	Execute(options map[string]bool) error
}