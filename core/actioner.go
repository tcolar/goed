package core

type Action interface {
	Run() error
}

type ActionDispatcher interface {
	Start()
	Dispatch(action Action)
	Shutdown()
}
