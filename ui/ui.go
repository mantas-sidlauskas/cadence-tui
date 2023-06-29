package ui

import (
	"github.com/jroimartin/gocui"
)

type Table interface {
	Render(g *gocui.Gui, v *gocui.View, maxX int)
	Header(v *gocui.View, maxX int)
	TableFromHere(selectedIndex int) (Table, error)
	Length() int
	SetGui(g *gocui.Gui)
	Name() string
}
