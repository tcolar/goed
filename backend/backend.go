// package backend provides the backend implementations of the goed
// editor text buffers.
package backend

import (
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/tcolar/goed/core"
)

func BufferFile(id int64) string {
	return path.Join(core.Home, "buffers", strconv.FormatInt(id, 10))
}

// NewMemBackendCmd creates a Comand runner backed by an In-memory based backend
// if title == nil then will show the command name
func NewMemBackendCmd(args []string, dir string, viewId int64, title *string) (*BackendCmd, error) {
	b, err := NewMemBackend("", viewId)
	if err != nil {
		return nil, err
	}
	c, err := newBackendCmd(args, dir, viewId, title)
	if err != nil {
		return nil, err
	}
	c.Backend = b
	c.Starter = &MemCmdStarter{}
	go c.Start(viewId)
	return c, nil
}

func newBackendCmd(args []string, dir string, viewId int64, title *string) (*BackendCmd, error) {
	c := &BackendCmd{
		dir:    dir,
		runner: exec.Command(args[0], args[1:]...),
		title:  title,
	}
	c.runner.Dir = dir
	if c.title == nil {
		title := strings.Join(c.runner.Args, " ")
		c.title = &title
	}
	return c, nil
}
