package backend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tcolar/goed/actions"
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
	go c.Start(c.ViewId())
	return nil
}

func (c *BackendCmd) Close() error {
	c.stop()
	return nil
}

func (c *BackendCmd) Running() bool {
	return c.runner != nil && c.runner.Process != nil
}

func (c *BackendCmd) Start(viewId int64) {
	workDir, _ := filepath.Abs(c.dir)
	actions.ViewSetWorkdirAction(viewId, workDir)
	c.runner.Dir = workDir
	actions.ViewSetTitleAction(viewId, fmt.Sprintf("[RUNNING] %s", *c.title))
	actions.ViewRenderAction(viewId)
	actions.EdTermFlushAction()

	err := c.Starter.Start(c)

	if err != nil {
		actions.ViewSetTitleAction(viewId, fmt.Sprintf("[FAILED] %s", *c.title))
		actions.EdSetStatusErrAction(err.Error())
	} else {
		actions.ViewSetTitleAction(viewId, *c.title)
	}
	actions.ViewSetWorkdirAction(viewId, workDir) // might have chnaged
	actions.EdRenderAction()
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
	w := backendAppender{backend: c.Backend, viewId: c.ViewId()}
	c.runner.Stdout = w
	c.runner.Stderr = w
	err := c.runner.Start()
	if err != nil {
		return err
	}
	err = c.runner.Wait()
	return err
}

type backendAppender struct {
	backend core.Backend
	viewId  int64
}

func (b backendAppender) Write(data []byte) (int, error) {
	var err error
	err = b.backend.Append(string(data))
	if err != nil {
		return 0, err
	}
	/*
				limit := core.Ed.Config().MaxCmdBufferLines
				if v == nil {
					return
				}
				if v.LineCount() > limit {
					c.Backend.Remove(1, 1, v.LineCount()-limit+1, 0)
				}

		event.ViewMoveCursorEvt(v, v.LineCount(), 0)
	*/
	//event.ViewMoveCursorEvt(v.LineCount(), 0)
	actions.ViewRenderAction(b.viewId)
	actions.EdTermFlushAction()
	return len(data), nil
}
