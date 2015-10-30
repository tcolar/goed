package backend

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/kr/pty"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

// BackendCmd is used to run a command using a specific backend
// whose content will be the the output of the command. (VT100 support)
//
// TODO : Deal with CTRL+C / Ctrl+R events, pass them through ?
// TODO : arrow events
type BackendCmd struct {
	core.Backend
	dir       string
	runner    *exec.Cmd
	pty       *os.File
	title     *string
	Starter   CmdStarter
	IsTerm    bool
	scrollTop bool // whether to scroll ack to top once command done
}

func (c *BackendCmd) Reload() error {
	if c.IsTerm {
		return errors.New("Can't reload terminal, close and reopen.")
	}
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

func (b *BackendCmd) Insert(row, col int, text string) error {
	if b.pty != nil {
		b.pty.Write([]byte(text))
	}
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
	c.pty.Write(data)
}

func (c *BackendCmd) Start(viewId int64) {
	workDir, _ := filepath.Abs(c.dir)
	actions.ViewSetWorkdir(viewId, workDir)
	c.runner.Dir = workDir
	actions.ViewSetTitle(viewId, fmt.Sprintf("[RUNNING] %s", *c.title))
	actions.ViewRender(viewId)
	actions.EdTermFlush()

	c.runner.Env = core.EnvWith([]string{"TERM=vt100",
		fmt.Sprintf("GOED_INSTANCE=%s", core.InstanceId),
		fmt.Sprintf("GOED_VIEW=%s", viewId)})

	err := c.Starter.Start(c)

	if err != nil {
		actions.EdSetStatusErr(err.Error())
		actions.ViewSetTitle(viewId, fmt.Sprintf("[FAILED] %s", *c.title))
	} else {
		actions.ViewSetTitle(viewId, *c.title)
	}
	actions.ViewSetWorkdir(viewId, workDir) // might have changed
	if c.scrollTop {
		actions.ViewCursorMvmt(viewId, core.CursorMvmtTop)
	}
	actions.EdRender()
}

func (c *BackendCmd) stop() {
	if c.runner != nil && c.runner.Process != nil {
		c.runner.Process.Kill()
		c.runner.Process = nil
	}
	if c.pty != nil {
		c.pty.Close()
		c.pty = nil
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
	t := core.Ed.Theme()
	w := backendAppender{backend: c.Backend, viewId: c.ViewId(), curFg: t.Fg, curBg: t.Bg}
	endc := make(chan struct{}, 1)
	go w.refresher(endc)
	var err error
	c.pty, err = pty.Start(c.runner)
	if err != nil {
		return err
	}

	go io.Copy(&w, c.pty)

	err = c.runner.Wait()
	close(endc)

	time.Sleep(100 * time.Millisecond)
	actions.EdRender()
	return err
}

type backendAppender struct {
	backend      core.Backend
	viewId       int64
	line, col    int
	dirty        int32      // >0 if dirty
	curFg, curBg core.Style // current terminal color attributes
}

// refresh the view if needed(dirty) but no more than every so often
// this greatly enhances performance and responsivness
func (b *backendAppender) refresher(endc chan struct{}) {
	for {
		select {
		case <-endc:
			actions.EdRender()
			return
		default:
			if atomic.SwapInt32(&b.dirty, 0) > 0 {
				actions.ViewTrim(b.viewId, core.Ed.Config().MaxCmdBufferLines)
				l, c := actions.ViewCurPos(b.viewId)
				actions.ViewMoveCursor(b.viewId, b.line-l, b.col-c)
				actions.EdRender()
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (b *backendAppender) Write(data []byte) (int, error) {
	size, err := b.vt100(data)
	atomic.AddInt32(&b.dirty, 1)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (b *backendAppender) flush(data []byte) error {
	m := b.backend.(*MemBackend)
	b.line, b.col = m.Overwrite(b.line, b.col, string(data), b.curFg, b.curBg)
	return nil
}
