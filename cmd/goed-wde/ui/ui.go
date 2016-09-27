// Goed using WDE as a graphical backend
// Experimental at this point, works but slow as high resolutions.
package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	wde "github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"github.com/tcolar/goed"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/event"
	"github.com/tcolar/goed/ui/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var _ core.Term = (*GuiTerm)(nil)

var palette = xtermPalette()

// backup/symbols font
var notoSymbols *truetype.Font

// it seem wde sometimes crash if we try to do things concurrently
// so using a lock
var wdeLock sync.Mutex

func Main() {
	config := goed.Initialize()
	if config == nil {
		return
	}
	term := NewGuiTerm(1200, 800, config)
	goed.Start(term, config)
}

// GuiTerm is a very minimal text terminal emulation GUI.
type GuiTerm struct {
	w, h             int
	text             [][]char
	win              wde.Window
	font             *truetype.Font
	charW, charH     int // size of characters
	face             font.Face
	ctx              *freetype.Context
	rgba             *image.RGBA
	cursorX, cursorY int
	fontPath         string
	fontSize         int
	fontDpi          int
}

type char struct {
	rune
	fg, bg        core.Style
	prevPaintHash uint64
}

func (c *char) hash() uint64 {
	return uint64(c.rune)<<32 | uint64(c.fg.Uint16())<<16 | uint64(c.bg.Uint16())
}

func NewGuiTerm(wh, ww int, config *core.Config) *GuiTerm {
	notoSymbols = parseBuiltinFont("fonts/NotoSansSymbols-Regular.ttf")

	t := &GuiTerm{
		fontPath: config.GuiFont,
		fontSize: config.GuiFontSize,
		fontDpi:  config.GuiFontDpi,
	}

	t.text = [][]char{}

	t.applyFont(t.fontPath, t.fontSize, wh, ww)

	return t
}

func (t *GuiTerm) applyFont(fontPath string, fontSize, wh, ww int) {
	if len(fontPath) == 0 {
		// builtin default font
		t.font = parseBuiltinFont("fonts/LiberationMono-Bold.ttf")
	} else {
		// user specified font
		t.font = parseFileFont(fontPath)
	}
	opts := truetype.Options{}
	opts.Size = float64(fontSize)
	t.face = truetype.NewFace(t.font, &opts)
	bounds, _, _ := t.face.GlyphBounds('░')
	t.charW = int((bounds.Max.X-bounds.Min.X)>>6) + t.fontDpi/32
	t.charH = int((bounds.Max.Y-bounds.Min.Y)>>6) + t.fontDpi/16

	t.ctx = freetype.NewContext()
	t.ctx.SetDPI(float64(t.fontDpi))
	t.ctx.SetFont(t.font)
	t.ctx.SetFontSize(float64(fontSize))
	t.ctx.SetHinting(font.HintingFull)

	t.resize(wh, ww)
}

func (t *GuiTerm) resize(ww, wh int) {
	w, h := t.w, t.h
	t.w = ww / t.charW
	t.h = wh / t.charH
	for i := 0; i < h; i++ {
		if t.w <= w {
			t.text[i] = t.text[i][:t.w] // truncate lines if needed
		} else {
			// expand lines if needed
			t.text[i] = append(t.text[i], make([]char, t.w-w)...)
		}
	}
	// extra lines if needed
	for i := h; i < t.h; i++ {
		t.text = append(t.text, make([]char, t.w))
	}
	// truncate number of lines if needed
	t.text = t.text[:t.h]
	// Update image/bounds
	t.rgba = image.NewRGBA(image.Rect(0, 0, ww, wh))
	t.ctx.SetClip(t.rgba.Bounds())
	t.ctx.SetDst(t.rgba)
	// force repaint all
	for y, ln := range t.text {
		for x, _ := range ln {
			t.text[y][x].prevPaintHash = 0
		}
	}
}

func parseBuiltinFont(fontPath string) *truetype.Font {
	fontBytes, err := fonts.Asset(fontPath)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		os.Exit(1)
	}
	return parseFont(fontBytes)
}

func parseFileFont(fontPath string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		os.Exit(1)
	}
	return parseFont(fontBytes)
}

func parseFont(fontBytes []byte) *truetype.Font {
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		os.Exit(1)
	}
	return font
}

func (t *GuiTerm) Init() error {
	return nil
}

func (t *GuiTerm) Close() {
	t.win.Close()
	wde.Stop()
}

func (t *GuiTerm) Clear(fg, bg core.Style) {
	zero := rune(0)
	for y, ln := range t.text {
		for x, _ := range ln {
			if t.text[y][x].rune != zero {
				t.text[y][x].rune = zero
			}
		}
	}
}

func (t *GuiTerm) Flush() {
	if t.win == nil {
		return
	}
	t.paint()
}

func (t *GuiTerm) SetCursor(x, y int) {
	if t.win == nil {
		return
	}
	px, py := t.cursorX, t.cursorY
	t.cursorX = x
	t.cursorY = y

	// force redraw where the cursor was/is
	if py >= 0 && py < len(t.text) && px >= 0 && px < len(t.text[py]) {
		t.text[py][px].prevPaintHash = 0
	}
	if py >= 0 && y < len(t.text) && x >= 0 && x < len(t.text[y]) {
		t.text[y][x].prevPaintHash = 0
	}
}

func (t *GuiTerm) Char(y, x int, c rune, fg, bg core.Style) {
	if x < 0 || y < 0 || y >= len(t.text) || x >= len(t.text[y]) {
		return
	}
	cur := &t.text[y][x]
	if cur.rune == c && cur.fg == fg && cur.bg == bg {
		return
	}

	cur.rune = c
	cur.fg = fg
	cur.bg = bg
}

// size in characters
func (t *GuiTerm) Size() (h, w int) {
	return t.h, t.w
}

// for testing
func (t *GuiTerm) CharAt(y, x int) rune {
	if x < 0 || y < 0 {
		panic("CharAt out of bounds")
	}
	if y >= t.h || x >= t.w {
		panic("CharAt out of bounds")
	}
	return t.text[y][x].rune
}

func (t *GuiTerm) SetExtendedColors(b bool) { // N/A
}

func (t *GuiTerm) Listen() {
	go t.listen()
	wde.Run()
}

func (t *GuiTerm) listen() {
	// For an unknow reason, wde won't bring the windows to the front if it's created
	// before this (OSX), so we wait until here to create it.
	win, err := wde.NewWindow(t.w*t.charW, t.h*t.charH)
	win.SetTitle("GoEd")
	if err != nil {
		panic(err)
	}
	t.win = win
	t.win.Show()
	// It seems a window maximize event BEFORE any paint causes a wde crash (OSX)
	// so making sure paint is called at least once
	t.paint()

	// Start the event loop
	evtState := event.NewEvent()
	dragY, dragX := 0, 0
	events := t.win.EventChan()
	for ev := range events {
		evtState.Type = event.Evt_None
		evtState.Glyph = ""
		evtState.MouseBtns[event.MouseWheelDown] = false // reset wheel down, no separate "up" event
		evtState.MouseBtns[event.MouseWheelUp] = false   // reset wheel up

		switch e := ev.(type) {
		case wde.ResizeEvent:
			t.resize(e.Width, e.Height)
			actions.Ar.EdResize(t.h, t.w)
			evtState.Type = event.EvtWinResize
		case wde.CloseEvent:
			evtState.Type = event.EvtQuit
			return
		case wde.MouseDownEvent:
			evtState.MouseDown(event.MouseButton(e.Which), e.Where.Y/t.charH, e.Where.X/t.charW)
		case wde.MouseUpEvent:
			evtState.MouseUp(event.MouseButton(e.Which), e.Where.Y/t.charH, e.Where.X/t.charW)
		case wde.ScrollEvent:
			mb := event.MouseWheelDown
			d := e.Delta
			switch {
			case d.Y > 0:
				mb = event.MouseWheelUp
			case d.X > 0:
				mb = event.MouseWheelRight
			case d.X < 0:
				mb = event.MouseWheelLeft
			}
			evtState.MouseDown(mb, e.Where.Y/t.charH, e.Where.X/t.charW)
		case wde.MouseDraggedEvent:
			// only send a drag event if moved to new text cell
			y, x := e.Where.Y/t.charH, e.Where.X/t.charW
			if y == dragY && x == dragX {
				continue
			}
			evtState.MouseDown(event.MouseButton(e.Which), y, x)
			dragX = x
			dragY = y
		case wde.KeyTypedEvent:
			evtState.Glyph = e.Glyph
		case wde.KeyDownEvent:
			evtState.KeyDown(e.Key)
			continue
		case wde.KeyUpEvent:
			evtState.KeyUp(e.Key)
			continue
		default:
			continue
		}
		event.Queue(*evtState)
	}
}

func (t *GuiTerm) paint() {
	rectangles := []image.Rectangle{}
	var curRect image.Rectangle
	any := false
	for y, ln := range t.text {
		for x, _ := range ln {
			if t.paintChar(y, x) { // char needs repainting
				// Try to reduce screen flush area to a few rectangles
				switch {
				case !any: // no rectangle yet
					curRect = image.Rect(x, y, x, y)
					any = true
				case curRect.Max.Y == y && x >= curRect.Min.X && x <= curRect.Max.X:
					//nada, already in the rectangle
				case y < curRect.Max.Y+3 && x > curRect.Min.X-3 && x < curRect.Max.X+3:
					curRect.Max.Y = y
					if x < curRect.Min.X {
						curRect.Min.X = x
					} else if x > curRect.Max.X {
						curRect.Max.X = x
					}
				default:
					// else start new rectangle
					rectangles = append(rectangles, curRect)
					curRect = image.Rect(x, y, x, y)
				}
			}
		}
	}
	if any {
		rectangles = append(rectangles, curRect)
	}
	// to UI coordinates
	for _, r := range rectangles {
		r.Min.X *= t.charW
		r.Max.X *= (t.charW + 1)
		r.Min.Y *= t.charH
		r.Max.Y *= (t.charH + 1)
	}
	if len(rectangles) > 0 {
		wdeLock.Lock()
		// unfortunately on OSX rectangles is not used and the whole screen is flushed
		t.win.FlushImage(rectangles...)
		wdeLock.Unlock()
	}
}

func (t *GuiTerm) paintChar(y, x int) bool {
	if y >= len(t.text) || x >= len(t.text[y]) {
		return false
	}
	atCursor := y == t.cursorY && x == t.cursorX

	r := &t.text[y][x]

	if r.rune < 32 {
		r.rune = ' '
		r.bg = core.Ed.Theme().Bg
	}

	hash := r.hash()
	if hash == r.prevPaintHash {
		return false
	}
	r.prevPaintHash = hash

	bg := image.NewUniform(palette[r.bg.Uint16()&255])
	fg := image.NewUniform(palette[r.fg.Uint16()&255])
	// cursor location gets inverted colors
	if atCursor {
		bg, fg = fg, bg
	}
	t.ctx.SetSrc(fg)
	rx := t.charW * x
	ry := t.charH * y
	rect := image.Rect(rx, ry, rx+t.charW, ry+t.charH)
	draw.Draw(t.rgba, rect, bg, image.ZP, draw.Src)
	pt := freetype.Pt(x*t.charW, t.charH-4+y*t.charH)
	t.drawRune(r.rune, pt)

	// copy it to wde
	wdeLock.Lock()
	i := t.rgba.SubImage(rect).(*image.RGBA)
	t.win.Screen().CopyRGBA(i, rect)
	wdeLock.Unlock()

	return true
}

// Draw the rune, if the user-picked font does not provide a glyph for the given
// rune try to fallback to notoSymbols
func (t *GuiTerm) drawRune(r rune, pt fixed.Point26_6) {
	if t.font.Index(r) != 0 {
		t.ctx.DrawString(string(r), pt)
		return
	}
	// if rune not found in main font, try symbol font
	font := t.font
	if notoSymbols.Index(r) != 0 {
		t.ctx.SetFont(notoSymbols)
		t.ctx.SetFontSize(float64(t.fontSize - 3))
	}
	t.ctx.DrawString(string(r), pt)
	t.ctx.SetFontSize(float64(t.fontSize))
	t.ctx.SetFont(font)
}

// Palette based of what's used in gnome-terminal / xterm-256
func xtermPalette() *[256]color.Color {
	a := uint8(255)
	// base colors (from gnome-terminal)
	palette := [256]color.Color{
		color.RGBA{0x2e, 0x34, 0x36, a},
		color.RGBA{0xcc, 0, 0, a},
		color.RGBA{0x4e, 0x9a, 0x06, a},
		color.RGBA{0xc4, 0xa0, 0, a},
		color.RGBA{0x34, 0x65, 0xa4, a},
		color.RGBA{0x75, 0x50, 0x7b, a},
		color.RGBA{0x06, 0x98, 0x9a, a},
		color.RGBA{0xd3, 0xd7, 0xcf, a},
		color.RGBA{0x55, 0x57, 0x53, a},
		color.RGBA{0xef, 0x29, 0x29, a},
		color.RGBA{0x8a, 0xe2, 0x34, a},
		color.RGBA{0xfc, 0xe9, 0x4f, a},
		color.RGBA{0x72, 0x9f, 0xcf, a},
		color.RGBA{0xad, 0x7f, 0xa8, a},
		color.RGBA{0x34, 0xe2, 0xe2, a},
		color.RGBA{0xee, 0xee, 0xec, a},
	}
	// xterm-256 colors
	for i := 16; i != 232; i++ {
		b := ((i - 16) % 6) * 40
		if b != 0 {
			b += 55
		}
		g := (((i - 16) / 6) % 6) * 40
		if g != 0 {
			g += 55
		}
		r := ((i - 16) / 36) * 40
		if r != 0 {
			r += 55
		}
		palette[i] = color.RGBA{uint8(r), uint8(g), uint8(b), a}
	}
	// Shades of grey
	for i := 232; i != 256; i++ {
		h := 8 + (i-232)*10
		palette[i] = color.RGBA{uint8(h), uint8(h), uint8(h), a}
	}

	return &palette
}
