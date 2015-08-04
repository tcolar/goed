package core

type Action interface {
	Run() error
}

type ActionDispatcher interface {
	Dispatch(action Action)
	Flush()
	Shutdown()
	Start()
}
