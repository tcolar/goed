package backend

import (
	"fmt"
	"io"
	"log"
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
type BackendCmd struct {
	*MemBackend
	dir       string
	runner    *exec.Cmd
	pty       *os.File
	title     *string
	Starter   CmdStarter
	scrollTop bool // whether to scroll back to top once command completed
	MaxRows   int  // ring buffer size
	head      int
}

func (c *BackendCmd) Reload() error {
	args, dir := c.runner.Args, c.runner.Dir
	c.stop()
	c.runner = exec.Command(args[0], args[1:]...)
	c.runner.Dir = dir
	c.MemBackend.Close()
	os.Remove(c.BufferLoc())
	c.MemBackend.Reload()
	go c.Start(c.ViewId())
	return nil
}

func (b *BackendCmd) ColorAt(ln, col int) (fg, bg core.Style) {
	if b.MaxRows > 0 {
		ln = (ln + b.head) % b.MaxRows
	}
	return b.MemBackend.ColorAt(ln, col)
}

func (b *BackendCmd) Insert(row, col int, text string) error {
	if b.pty != nil {
		b.pty.Write([]byte(text))
	}
	return nil
}

func (b *BackendCmd) Append(text string) error {
	return fmt.Errorf("Not implemented, Append()")
}

func (b *BackendCmd) Remove(row1, col1, row2, col2 int) error {
	return fmt.Errorf("Not implemented, Insert()")
}

func (b *BackendCmd) Save(loc string) error {
	return fmt.Errorf("Not implemented, Save()")
}

func (b *BackendCmd) SrcLoc() string {
	return ""
}

func (b *BackendCmd) BufferLoc() string {
	return "_MEM_"
}

func (b *BackendCmd) Close() error {
	b.stop()
	b.MemBackend.Close()
	return nil
}

func (b *BackendCmd) Running() bool {
	return b.runner != nil && b.runner.Process != nil
}

func (b *BackendCmd) SendBytes(data []byte) {
	b.pty.Write(data)
}

func (b *BackendCmd) Head() int { // for unit testing
	return b.head
}

func (c *BackendCmd) Start(viewId int64) {
	workDir, _ := filepath.Abs(c.dir)
	actions.ViewSetWorkdir(viewId, workDir)
	c.runner.Dir = workDir
	actions.ViewSetTitle(viewId, fmt.Sprintf("[RUNNING] %s", *c.title))
	actions.ViewRender(viewId)
	actions.EdTermFlush()

	c.runner.Env = core.EnvWith([]string{"TERM=vt100",
		fmt.Sprintf("GOED_INSTANCE=%d", core.InstanceId),
		fmt.Sprintf("GOED_VIEW=%d", viewId)})

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

func (b *BackendCmd) Wipe() {
	b.MemBackend.lock.Lock()
	b.head = 0
	b.MemBackend.lock.Unlock()
	b.MemBackend.Wipe()
}

func (b *BackendCmd) adjustHead(row int) {
	h := core.Ed.Config().SyntaxHighlighting
	head := (b.head + row - b.MaxRows + 1) % b.MaxRows
	if head != b.head {
		// if head changed, we are going to be reusing a line
		// clear it on first re-use.
		for i := b.head; i != head; i++ {
			if i >= b.MaxRows {
				i = -1
				continue
			}
			if i >= len(b.text) {
				continue
			}
			b.text[i] = b.text[i][:0]
			if h {
				b.colors[i] = b.colors[i][:0]
			}
		}
		b.head = head
	}
}

func (b *BackendCmd) Overwrite(row, col int, text string, fg, bg core.Style) (atRow, atCol int) {
	if len(text) == 0 {
		return row, col
	}

	h := core.Ed.Config().SyntaxHighlighting

	b.MemBackend.lock.Lock()
	defer b.MemBackend.lock.Unlock()

	runes := core.StringToRunes(text)
	for z, ln := range runes { // each line
		if z > 0 {
			row++
			col = 0
		}
		r := row
		if b.MaxRows > 0 {
			r = (b.head + row) % b.MaxRows
			if row >= b.MaxRows {
				b.adjustHead(row)
			}
		}
		for _, ch := range ln {
			if col >= b.MemBackend.vtCols { // wrap lines wider than terminal width
				col = 0
				row++
				if b.MaxRows > 0 {
					r = (b.head + row) % b.MaxRows
					if row >= b.MaxRows {
						b.adjustHead(row)
					}
				}
			}
			if len(b.text) <= r { // make sure we have enough rows
				b.text = append(b.text, make([][]rune, r-len(b.text)+1)...)
			}
			if h && len(b.colors) <= r {
				b.colors = append(b.colors, make([][]*color, r-len(b.colors)+1)...)
			}
			if len(b.text[r]) <= col { // make sure enough cols
				b.text[r] = append(b.text[r], make([]rune, col-len(b.text[r])+1)...)
			}
			if h && len(b.colors[r]) <= col {
				b.colors[r] = append(b.colors[r], make([]*color, col-len(b.colors[r])+1)...)
			}
			b.text[r][col] = ch // write the char
			if h {
				b.colors[r][col] = &color{fg, bg}
			}
			col++
		}
	}
	if b.MaxRows > 0 && row >= b.MaxRows {
		return b.MaxRows - 1, col
	}
	return row, col
}

func (b *BackendCmd) Slice(row, col, row2, col2 int) *core.Slice {
	if b.MaxRows <= 0 || b.head == 0 {
		return b.MemBackend.Slice(row, col, row2, col2)
	}
	// Otherwise ringbuffer
	b.MemBackend.lock.Lock()
	defer b.MemBackend.lock.Unlock()
	slice := core.NewSlice(row, col, row2, col2, [][]rune{})
	text := slice.Text()

	r := (row + b.head) % b.MaxRows
	r2 := (row2 + b.head) % b.MaxRows
	if row2 == -1 || row2 >= len(b.text) {
		// then read whole buffer (wrap all the way around)
		r2 = (b.head - 1) % b.MaxRows
	}

	if r < 0 || r2 < 0 {
		return slice
	}

	if r2 >= r {
		// no wrap-around
		*text = append(*text, *b.MemBackend.sliceNoLock(r, col, r2, col2).Text()...)
		return slice
	}
	// otherwise read end + beginning
	*text = append(*text, *b.MemBackend.sliceNoLock(r, col, -1, col2).Text()...)
	*text = append(*text, *b.MemBackend.sliceNoLock(0, col, r2, col2).Text()...)
	return slice
}

func (b *BackendCmd) clearLn(row, col int) {
	h := core.Ed.Config().SyntaxHighlighting
	if b.MaxRows > 0 {
		row = (row + b.head) % b.MaxRows
	}
	b.MemBackend.lock.Lock()
	defer b.MemBackend.lock.Unlock()
	if row >= len(b.text) || col >= len(b.text[row]) {
		return
	}
	b.text[row] = b.text[row][:col]
	if h {
		b.colors[row] = b.colors[row][:col]
	}
}

func (b *BackendCmd) clearScreen(row, col int) {
	h := core.Ed.Config().SyntaxHighlighting
	b.lock.Lock()
	defer b.lock.Unlock()
	if row >= len(b.text) {
		return
	}
	if b.MaxRows <= 0 {
		b.text = b.text[0 : row+1]
		if h {
			b.colors = b.colors[0 : row+1]
		}
		if col < len(b.text[row]) {
			b.text[row] = b.text[row][:col]
			if h {
				b.colors[row] = b.colors[row][:col]
			}
		}
		return
	}
	// else ringbuffer
	r := (b.head + row) % len(b.text)
	b.text[r] = b.text[r][:col]
	if h {
		b.colors[r] = b.colors[r][:col]
	}
	r = (r + 1) % len(b.text)
	for r != b.head {
		b.text[r] = b.text[r][:0]
		if h {
			b.colors[r] = b.colors[r][:0]
		}
		r = (r + 1) % len(b.text)
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

	c.Wipe()
	return c.stream()
}

func (c *BackendCmd) stream() error {
	t := core.Ed.Theme()
	w := backendAppender{backend: c, viewId: c.ViewId(), curFg: t.Fg, curBg: t.Bg}
	endc := make(chan struct{}, 1)
	go w.refresher(endc)
	var err error
	c.pty, err = pty.Start(c.runner)
	if err != nil {
		return err
	}

	go func() {
		io.Copy(&w, c.pty)
		log.Println("copy done")
	}()

	err = c.runner.Wait()
	close(endc)

	time.Sleep(100 * time.Millisecond)
	actions.EdRender()
	return err
}

type backendAppender struct {
	backend      *BackendCmd
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
				// In case of view size change, adjust tty setting and wrap value.
				// TODO : No easy way to do that if the terminal is currently
				// in an interactive program
				//if v != nil && (rows != v.LastViewLine() || cols != v.LastViewCol()) {
				// refresh view
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
	b.line, b.col = b.backend.Overwrite(b.line, b.col, string(data), b.curFg, b.curBg)
	return nil
}
