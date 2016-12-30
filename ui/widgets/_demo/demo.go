// -build
package main

import (
	"time"

	"github.com/tcolar/goed/cmd/goed-tcell/ui"
	"github.com/tcolar/goed/ui/widgets"
)

func main() {
	term := ui.NewTcell()
	term.Init()
	term.SetExtendedColors(true)
	defer term.Close()

	w := widgets.NewTermWidget(term)
	q := widgets.NewDialog(
		"Continue ?\nwhat da ya think ??",
		0,
		widgets.NewButton("Ok", 'O', nil),
		widgets.NewButton("Maybe", 'M', nil),
		widgets.NewButton("No way", 'N', nil),
	)
	q.Render()
	w.AddWidget(q)
	w.Render()
	time.Sleep(5 * time.Second)
}
