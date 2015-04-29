package backend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tcolar/goed/core"
)

// BackendCmd is used to run a command using a specific backend
// whose content is the output from an extenal command.
type BackendCmd struct {
	core.Backend
	dir     string
	runner  *exec.Cmd
	title   *string
	Starter CmdStarter
}

func (c *BackendCmd) Reload() error {
	args, dir := c.runner.Args, c.runner.Dir
	c.stop()
	// It does not seem we can reuse a command so create a new one
	c.runner = exec.Command(args[0], args[1:]...)
	c.runner.Dir = dir
	c.Backend.Close()
	os.Remove(c.BufferLoc())
	c.Backend.Reload()
	go c.Start()
	return nil
}

func (c *BackendCmd) Close() error {
	c.stop()
	c.runner = nil
	return nil
}

func (c *BackendCmd) Start() {
	workDir, _ := filepath.Abs(c.dir)
	v := core.Ed.ViewById(c.ViewId())
	v.SetWorkDir(workDir)
	c.runner.Dir = workDir
	v.SetTitle(fmt.Sprintf("[RUNNING] %s", *c.title))
	core.Ed.Render()

	err := c.Starter.Start(c)

	if err != nil {
		v.SetTitle(fmt.Sprintf("[FAILED] %s", *c.title))
		core.Ed.SetStatusErr(err.Error())
	} else {
		v.SetTitle(*c.title)
	}
	v.SetWorkDir(workDir) // start() could have modified this
	core.Ed.Render()
}

func (c *BackendCmd) stop() {
	if c.runner != nil && c.runner.Process != nil {
		c.runner.Process.Release()
		c.runner.Process.Kill()
		c.runner.Process = nil
	}
}

type CmdStarter interface {
	Start(c *BackendCmd) error
}

// starter impl for file backend
type FileCmdStarter struct {
}

func (s *FileCmdStarter) Start(c *BackendCmd) error {
	b := c.Backend.(*FileBackend)
	c.runner.Stdout = b.file
	c.runner.Stderr = b.file
	err := c.runner.Run()
	b.Reload()
	return err
}

// starter impl for mem backend
type MemCmdStarter struct {
}

func (s *MemCmdStarter) Start(c *BackendCmd) error {

	b := c.Backend.(*MemBackend)
	out, err := c.runner.CombinedOutput()
	b.Wipe()
	b.Insert(1, 1, string(out))
	return err
}
