package backend

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tcolar/goed/core"
)

// BackendCmd is used to run a command using a specific backend
// whose content will be the the output of the command.
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
	return nil
}

func (c *BackendCmd) Running() bool {
	return c.runner != nil && c.runner.Process != nil
}

func (c *BackendCmd) Start() {
	workDir, _ := filepath.Abs(c.dir)
	v := core.Ed.ViewById(c.ViewId())
	v.SetWorkDir(workDir)
	c.runner.Dir = workDir
	v.SetTitle(fmt.Sprintf("[RUNNING] %s", *c.title))
	v.Render()
	core.Ed.TermFlush()

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
		c.runner.Process.Kill()
		c.runner.Process = nil
	}
}

// CmdStarter is an interface for a "startable" command
type CmdStarter interface {
	Start(c *BackendCmd) error
}

/*
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
}*/

// MemCmdStarter is the command starter implementation for mem backend
// It starts the command and "streams" the output to the backend.
type MemCmdStarter struct {
}

func (s *MemCmdStarter) Start(c *BackendCmd) error {

	b := c.Backend.(*MemBackend)
	b.Wipe()
	return c.stream()
}

func (c *BackendCmd) stream() error {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	c.runner.Stdout = &outBuf
	c.runner.Stderr = &errBuf
	err := c.runner.Start()
	if err != nil {
		return err
	}
	done := false
	go func() {
		buf := make([]byte, 50000)
		for {
			c.flush(&outBuf, &errBuf, buf)
			if done {
				c.flush(&outBuf, &errBuf, buf)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	err = c.runner.Wait()
	done = true
	return err
}

func (c *BackendCmd) flush(o, e *bytes.Buffer, buf []byte) {
	refresh := false
	v := core.Ed.ViewById(c.ViewId())
	if o.Len() > 0 {
		nb, _ := o.Read(buf)
		if nb > 0 {
			c.Backend.Append(string(buf[:nb]))
			refresh = true
		}
	}
	if e.Len() > 0 {
		nb, _ := e.Read(buf)
		if nb > 0 {
			c.Backend.Append(string(buf[:nb]))
			refresh = true
		}
	}
	if refresh {
		limit := core.Ed.Config().MaxCmdBufferLines
		if v.LineCount() > limit {
			c.Backend.Remove(1, 1, v.LineCount()-limit+1, 0)
		}
		v.MoveCursor(0, v.LineCount())
		v.Render()
		core.Ed.TermFlush()
	}
}
