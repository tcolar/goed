package core

type FileOp uint32

const (
	OpCreate FileOp = 1 << iota
	OpWrite
	OpRemove
	OpRename
	OpChmod
)
