package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	wde "github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/event"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var _ core.Term = (*GuiTerm)(nil)

var palette = xtermPalette()

// TODO: font config
var fontPath = "fonts/LiberationMono-Regular.ttf"

// backup fonts
var noto *truetype.Font
var notoSymbols *truetype.Font

var fontSize = 10
var dpi = 96

// GuiTerm is a very minimal text terminal emulation GUI.
type GuiTerm struct {
	w, h int
	text [][]char
	//	textLock         sync.RWMutex
	win              wde.Window
	font             *truetype.Font
	charW, charH     int // size of characters
	face             font.Face
	ctx              *freetype.Context
	rgba             *image.RGBA
	cursorX, cursorY int
}

type char struct {
	rune
	fg, bg core.Style
	fresh  bool
}

func (c char) equals(c2 char) bool {
	if c.rune != c2.rune {
		return false
	}
	if c.fg != c2.fg {
		return false
	}
	if c.bg != c2.bg {
		return false
	}
	return true
}

func NewGuiTerm(h, w int) *GuiTerm {
	noto = parseFont("fonts/NotoSans-Regular.ttf")
	notoSymbols = parseFont("fonts/NotoSansSymbols-Regular.ttf")

	win, err := wde.NewWindow(h, w)
	win.SetTitle("GoEd")
	if err != nil {
		panic(err)
	}

	t := &GuiTerm{
		win: win,
	}

	t.text = [][]char{}

	t.applyFont(fontPath, fontSize)

	return t
}

func (t *GuiTerm) applyFont(fontPath string, fontSize int) {
	t.font = parseFont(fontPath)
	opts := truetype.Options{}
	opts.Size = float64(fontSize)
	t.face = truetype.NewFace(t.font, &opts)
	bounds, _, _ := t.face.GlyphBounds('░')
	t.charW = int((bounds.Max.X-bounds.Min.X)>>6) + dpi/32
	t.charH = int((bounds.Max.Y-bounds.Min.Y)>>6) + dpi/16

	t.ctx = freetype.NewContext()
	t.ctx.SetDPI(float64(dpi))
	t.ctx.SetFont(t.font)
	t.ctx.SetFontSize(float64(fontSize))
	t.ctx.SetHinting(font.HintingFull)

	t.resize(t.win.Size())
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
	//fmt.Printf("%v %v %v\n", t.w, t.h, t.rgba.Bounds())
}

func parseFont(fontPath string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		os.Exit(1)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		os.Exit(1)
	}
	return font
}

func (t *GuiTerm) Init() error {
	t.win.Show()
	return nil
}

func (t *GuiTerm) Close() {
	t.win.Close()
	wde.Stop()
}

func (t *GuiTerm) Clear(fg, bg uint16) {
	zero := rune(0)
	for y, ln := range t.text {
		for x, _ := range ln {
			if t.text[y][x].rune != zero {
				t.text[y][x].rune = zero
				t.text[y][x].fresh = false
			}
		}
	}
}

func (t *GuiTerm) Flush() {
	t.paint()
}

func (t *GuiTerm) SetCursor(x, y int) {
	// todo : move cursor
	px, py := t.cursorX, t.cursorY
	t.cursorX = x
	t.cursorY = y

	t.paintChar(py, px)
	t.paintChar(y, x)

	ny, nx := y*t.charH, x*t.charW
	r := image.Rect(nx, ny, nx+t.charW, ny+t.charH)
	i := t.rgba.SubImage(r).(*image.RGBA)
	t.win.Screen().CopyRGBA(i, r)
	t.win.FlushImage()
}

func (t *GuiTerm) Char(y, x int, c rune, fg, bg core.Style) {
	if x < 0 || y < 0 || y >= len(t.text) || x >= len(t.text[y]) {
		return
	}
	ch := char{
		rune: c,
		fg:   fg,
		bg:   bg,
	}
	if !ch.equals(t.text[y][x]) {
		t.text[y][x] = ch
	}
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
	evtState := event.NewEvent()
	dragY, dragX := 0, 0
	for ev := range t.win.EventChan() {
		evtState.Type = event.Evt_None
		evtState.Glyph = ""
		switch e := ev.(type) {
		case wde.ResizeEvent:
			t.resize(e.Width, e.Height)
			actions.Ar.EdResize(t.h, t.w)
			evtState.Type = event.EvtWinResize
		case wde.CloseEvent:
			evtState.Type = event.EvtQuit
			return
		case wde.MouseDownEvent:
			evtState.MouseDown(int(e.Which), e.Where.Y/t.charH, e.Where.X/t.charW)
		case wde.MouseUpEvent:
			evtState.MouseUp(int(e.Which), e.Where.Y/t.charH, e.Where.X/t.charW)
		case wde.MouseDraggedEvent:
			// only send drag event if moved to new text cell
			y, x := e.Where.Y/t.charH, e.Where.X/t.charW
			if y == dragY && x == dragX {
				continue
			}
			evtState.MouseDown(int(e.Which), y, x)
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
	start := time.Now()
	for y, ln := range t.text {
		for x, _ := range ln {
			t.paintChar(y, x)
		}
	}
	fmt.Printf("paint1 %v\n", time.Now().Sub(start))
	t.win.Screen().CopyRGBA(t.rgba, t.rgba.Bounds())
	t.win.FlushImage()
	fmt.Printf("paint2 %v\n", time.Now().Sub(start))

}

func (t *GuiTerm) paintChar(y, x int) {
	if y >= len(t.text) || x >= len(t.text[y]) {
		return
	}
	r := t.text[y][x]
	//if r.fresh {
	//	return
	//}
	//t.text[y][x].fresh = true
	//fmt.Printf("%d,%d -> %s\n", y, x, string(r.rune))

	pt := freetype.Pt(1+x*t.charW, t.charH-4+y*t.charH)
	if r.rune < 32 {
		r.rune = ' '
		r.bg = core.Ed.Theme().Bg
	}
	// TODO: attributes (bold)
	bg := image.NewUniform(palette[r.bg.Uint16()&255])
	fg := image.NewUniform(palette[r.fg.Uint16()&255])
	// cursor location gets inverted colors
	if y == t.cursorY && x == t.cursorX {
		bg, fg = fg, bg
	}
	t.ctx.SetSrc(fg)
	rx := t.charW * x
	ry := t.charH * y
	rect := image.Rect(rx, ry, rx+t.charW, ry+t.charH)
	draw.Draw(t.rgba, rect, bg, image.ZP, draw.Src)
	t.drawRune(r.rune, pt)
}

// Draw the rune, if the user-picked font does not provide a glyph for the given
// rune try to fallback to noto / notoSymbols
func (t *GuiTerm) drawRune(r rune, pt fixed.Point26_6) {
	if t.font.Index(r) != 0 {
		t.ctx.DrawString(string(r), pt)
		return
	}
	font := t.font
	if font.Index(r) != 0 {
	} else if noto.Index(r) != 0 {
		t.ctx.SetFont(noto)
		t.ctx.SetFontSize(float64(fontSize - 3)) // "fat" font

	} else if notoSymbols.Index(r) != 0 {
		t.ctx.SetFont(notoSymbols)
		t.ctx.SetFontSize(float64(fontSize - 3))
	}
	t.ctx.DrawString(string(r), pt)
	// restore font
	t.ctx.SetFontSize(float64(fontSize))
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
