package backend

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

// BackendCmd is used to run a command using a specific backend
// whose content will be the the output of the command.
type BackendCmd struct {
	core.Backend
	dir     string
	runner  *exec.Cmd
	inPipe  io.WriteCloser
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

func (c *BackendCmd) SendBytes(data []byte) {
	c.inPipe.Write(data)
}

func (c *BackendCmd) Start(viewId int64) {
	workDir, _ := filepath.Abs(c.dir)
	actions.ViewSetWorkdir(viewId, workDir)
	c.runner.Dir = workDir
	actions.ViewSetTitle(viewId, fmt.Sprintf("[RUNNING] %s", *c.title))
	actions.ViewRender(viewId)
	actions.EdTermFlush()

	err := c.Starter.Start(c)

	if err != nil {
		actions.ViewSetTitle(viewId, fmt.Sprintf("[FAILED] %s", *c.title))
		actions.EdSetStatusErr(err.Error())
	} else {
		actions.ViewSetTitle(viewId, *c.title)
	}
	actions.ViewSetWorkdir(viewId, workDir) // might have changed
	actions.EdRender()
}

func (c *BackendCmd) stop() {
	if c.inPipe != nil {
		c.inPipe.Close()
	}
	if c.runner != nil && c.runner.Process != nil {
		c.runner.Process.Kill()
		c.runner.Process = nil
	}
}

// CmdStarter is an interface for a "startable" command
type CmdStarter interface {
	Start(c *BackendCmd) error
}

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
	c.inPipe, _ = c.runner.StdinPipe()
	err := c.runner.Start()
	if err != nil {
		return err
	}
	err = c.runner.Wait()

	actions.EdRender()
	return err
}

type backendAppender struct {
	backend   core.Backend
	viewId    int64
	lastFlush int64
}

func (b backendAppender) Write(data []byte) (int, error) {
	err := b.backend.Append(string(data))
	if err != nil {
		return 0, err
	}
	actions.ViewTrim(b.viewId, core.Ed.Config().MaxCmdBufferLines)

	actions.ViewCursorMvmt(b.viewId, core.CursorMvmtBottom)

	now := time.Now().Unix()

	// render every so often
	if now > b.lastFlush+500 {
		b.lastFlush = now
		actions.EdRender()
	}
	return len(data), nil
}
