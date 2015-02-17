package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BackendCmd is used to run a command using a specific backend
// whose content is the output from an extenal command.
type BackendCmd struct {
	Backend
	dir     string
	runner  *exec.Cmd
	title   *string
	starter cmdStarter
}

// Cmd runner with File based backend
// if title == nil then will show the command name
func (e *Editor) NewFileBackendCmd(args []string, dir string, viewId int, title *string) (*BackendCmd, error) {
	loc := e.BufferFile(viewId)
	os.Remove(loc)
	b, err := e.NewFileBackend(loc, viewId)
	if err != nil {
		return nil, err
	}
	c, err := e.newBackendCmd(args, dir, viewId, title)
	if err != nil {
		return nil, err
	}
	c.Backend = b
	c.starter = &fileCmdStarter{}
	go c.start()
	return c, nil
}

// Cmd runner with In-memory based backend
// if title == nil then will show the command name
func (e *Editor) NewMemBackendCmd(args []string, dir string, viewId int, title *string) (*BackendCmd, error) {
	b, err := e.NewMemBackend("", viewId)
	if err != nil {
		return nil, err
	}
	c, err := e.newBackendCmd(args, dir, viewId, title)
	if err != nil {
		return nil, err
	}
	c.Backend = b
	c.starter = &memCmdStarter{}
	go c.start()
	return c, nil
}

func (e *Editor) newBackendCmd(args []string, dir string, viewId int, title *string) (*BackendCmd, error) {
	c := &BackendCmd{
		dir:    dir,
		runner: exec.Command(args[0], args[1:]...),
		title:  title,
	}
	if c.title == nil {
		title := strings.Join(c.runner.Args, " ")
		c.title = &title
	}
	return c, nil
}

func (c *BackendCmd) ReRun() {
	c.stop()
	go c.start()
}

func (c *BackendCmd) Close() error {
	c.stop()
	c.runner = nil
	return nil
}

func (c *BackendCmd) start() {
	workDir, _ := filepath.Abs(c.dir)
	v := Ed.ViewById(c.ViewId())
	v.WorkDir = workDir
	c.runner.Dir = workDir
	v.title = fmt.Sprintf("[RUNNING] %s", *c.title)
	Ed.Render()

	err := c.starter.start(c, v)

	if err != nil {
		v.title = fmt.Sprintf("[FAILED] %s", *c.title)
		Ed.SetStatusErr(err.Error())
	} else {
		v.title = *c.title
	}
	v.WorkDir = workDir // start() could have modified this
	Ed.Render()
}

func (c *BackendCmd) stop() {
	if c.runner != nil && c.runner.Process != nil {
		c.runner.Process.Release()
		c.runner.Process.Kill()
	}
}

type cmdStarter interface {
	start(c *BackendCmd, v *View) error
}

// starter impl for file backend
type fileCmdStarter struct {
}

func (s *fileCmdStarter) start(c *BackendCmd, v *View) error {
	b := c.Backend.(*FileBackend)
	c.runner.Stdout = b.file
	c.runner.Stderr = b.file
	err := c.runner.Run()
	Ed.Open(c.SrcLoc(), v, "")
	return err
}

// starter impl for mem backend
type memCmdStarter struct {
}

func (s *memCmdStarter) start(c *BackendCmd, v *View) error {

	b := c.Backend.(*MemBackend)
	out, err := c.runner.CombinedOutput()
	// TODO: clear ??
	b.Insert(1, 1, string(out))
	return err
}
