package core

type Action interface {
	Run()
}

type ActionDispatcher interface {
	Dispatch(action Action)
	Flush()
	Shutdown()
	Start()
}
