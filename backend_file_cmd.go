package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileBackendCmd represents a buffer file,
// whose content is the output from an extenal command.
type FileBackendCmd struct {
	FileBackend
	dir    string
	runner *exec.Cmd
	view   *View
}

func (e *Editor) NewFileBackendCmd(args []string, dir string, v *View) (*FileBackendCmd, error) {
	b, err := e.NewFileBackend(e.BufferFile(v.Id), v)
	if err != nil {
		return nil, err
	}
	fb := &FileBackendCmd{
		FileBackend: *b,
		dir:         dir,
		runner:      exec.Command(args[0], args[1:]...),
	}
	// TODO: go fb.start()
	return fb, nil
}

func (f *FileBackendCmd) Refresh() {
	f.stop()
	go f.start()
}

func (f *FileBackendCmd) Close() error {
	f.stop()
	return nil
}

func (f *FileBackendCmd) start() {
	workDir, _ := filepath.Abs(f.dir)
	v := f.view
	v.WorkDir = workDir
	f.runner.Stdout = f.file
	f.runner.Stderr = f.file
	f.runner.Dir = workDir
	title := strings.Join(f.runner.Args, " ")
	v.title = fmt.Sprintf("[RUNNING] %s", title)
	Ed.Render()
	err := f.runner.Run()
	Ed.Open(f.srcLoc, v, "")
	v.WorkDir = workDir // open() would have modified this
	if err != nil {
		v.title = fmt.Sprintf("[FAILED] %s", title)
		Ed.SetStatusErr(err.Error())
	} else {
		v.title = fmt.Sprintf("[DONE] %s", title)
		Ed.SetStatus(workDir)
	}
	Ed.Render()
}

func (f *FileBackendCmd) stop() {
	// TODO: kill command properly
	f.runner = nil
}

/*
// FileBackendInternalCmd is a special type of Command where the command is "internal" to goed,
// rather than an exeternal command.
type FileBackendInternalCmd struct {
	FileBackend
	cmd string
	dir string
}

type CmdName string

const (
	CmdDirLs CmdName = "ls"
)

func NewFileBackendInternalCmd(cmd CmdName, dir string, bufferId int) (*FileBackendCmd, error) {
	b, err := NewFileBackend("TODO", bufferId)
	if err != nil {
		return nil, err
	}
	return &FileBackendCmd{
		FileBackend: *b,
		cmd:         cmd,
		dir:         dir,
	}, nil
}

func (f *FileBackendInternalCmd) Refresh() {
	//TODO
}*/
