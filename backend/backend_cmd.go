package backend

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kr/pty"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

// TODO : Handle VT100 codes
// TODO : Handle color codes (or ignore if no colors)
// TODO : Deal with CTRL+C / Ctrl+R events, pass them through ?
// TODO : arrow events ?
// TODO : custom "cd" command so it chnages view workdir too

// BackendCmd is used to run a command using a specific backend
// whose content will be the the output of the command.
type BackendCmd struct {
	core.Backend
	dir     string
	runner  *exec.Cmd
	pty     *os.File
	title   *string
	Starter CmdStarter
	IsTerm  bool
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
	w := backendAppender{backend: c.Backend, viewId: c.ViewId()}
	var err error
	c.pty, err = pty.Start(c.runner)
	if err != nil {
		return err
	}

	go io.Copy(w, c.pty)

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
	size, err := b.vt100(data)
	if err != nil {
		return 0, err
	}

	actions.ViewTrim(b.viewId, core.Ed.Config().MaxCmdBufferLines)

	actions.ViewCursorMvmt(b.viewId, core.CursorMvmtBottom)

	now := time.Now().Unix()

	// render at most every so often
	if now > b.lastFlush+100 {
		b.lastFlush = now
		actions.EdRender()
	}
	return size, nil
}

func (b backendAppender) vt100(data []byte) (int, error) {
	from := 0
	for i := 0; i < len(data); i++ {
		if b.consumeVt100(data, from, &i) {
			from = i
			i--
		}
	}
	// flush leftover
	err := b.flush(data[from:len(data)])
	return len(data), err
}

func (b backendAppender) consumeVt100(data []byte, from int, i *int) bool {
	start := *i
	// set title
	if b.consume(data, i, 27) && b.consume(data, i, 93) &&
		b.consume(data, i, '0') && b.consume(data, i, ';') {
		b.consumeUntil(data, i, 7)
		actions.ViewSetTitle(b.viewId, string(data[start+4:*i]))
		return true
	}
	*i = start
	// Start real VT100 codes
	if !(b.consume(data, i, 27) && b.consume(data, i, 91)) { // ^[
		return false
	}
	*i = start + 2
	// Color attribute + fg color  + bg color
	if b.consumeNumber(data, i) && b.consume(data, i, ';') &&
		b.consumeNumber(data, i) && b.consume(data, i, ';') &&
		b.consumeNumber(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start]) // flush what was before escape seq first.
		return true
	}
	// Color attribute + fg color
	*i = start + 2
	if b.consumeNumber(data, i) && b.consume(data, i, ';') &&
		b.consumeNumber(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start])
		return true
	}
	*i = start + 2
	if b.consumeNumber(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start])
		return true
	}
	// Various Set comand, ignore for now
	*i = start + 2
	if b.consume(data, i, '?') && b.consumeNumber(data, i) && b.consume(data, i, 'h') {
		b.flush(data[from:start])
		return true
	}
	// Various Set comand, ignore for now
	*i = start + 2
	if b.consume(data, i, '?') && b.consumeNumber(data, i) && b.consume(data, i, 'l') {
		b.flush(data[from:start])
		return true
	}
	*i = start + 2
	// Set alternate keypad mode
	if b.consume(data, i, '=') {
		b.flush(data[from:start])
		return true
	}
	*i = start + 2
	// Set numeric keypad mode
	if b.consume(data, i, '>') {
		b.flush(data[from:start])
		return true
	}
	*i = start + 2
	// clear to EOL
	if b.consume(data, i, 'J') {
		b.flush(data[from:start])
		return true
	}
	*i = start + 2
	// clear to EOF
	if b.consume(data, i, 'K') {
		b.flush(data[from:start])
		return true
	}
	// no match
	*i = start
	return false
}

func (b backendAppender) consumeNumber(data []byte, i *int) bool {
	found := false
	for *i < len(data) && data[*i] >= '0' && data[*i] <= '9' {
		*i++
		found = true
	}
	return found
}

func (b backendAppender) consume(data []byte, i *int, c byte) bool {
	if *i >= len(data) {
		return false
	}
	if data[*i] == c {
		*i++
		return true
	}
	return false
}

func (b backendAppender) consumeUntil(data []byte, i *int, c byte) {
	for *i < len(data) && data[*i] != c {
		*i++
	}
	return
}

func (b backendAppender) flush(data []byte) error {
	return b.backend.Append(string(data))
}
