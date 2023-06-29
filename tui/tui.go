package tui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/mantas-sidlauskas/cadence-tui/cadence"
	"github.com/mantas-sidlauskas/cadence-tui/config"
	"github.com/mantas-sidlauskas/cadence-tui/ui"
	"go.uber.org/multierr"
)

type TUI struct {
	g       *gocui.Gui
	config  *config.Config
	current ui.Table

	previousTables  []ui.Table
	previousCursors []int
	previousOrigins []int
	cc              *cadence.Client
}

func New(config *config.Config) *TUI {
	return &TUI{
		config: config,
		cc:     cadence.NewClient(config.Address),
	}
}

func (t *TUI) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	g.InputEsc = true
	if err != nil {
		return err
	}
	defer g.Close()
	t.g = g

	t.current = ui.NewDomainsTable(t.cc)
	t.g.SetManagerFunc(t.layout)

	if err := t.keybindings(); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil

}

func (t *TUI) keybindings() error {
	err := multierr.Combine(
		t.g.SetKeybinding("", 'q', gocui.ModNone, t.quit),
		t.g.SetKeybinding("table", 's', gocui.ModNone, t.showClusterInfo),
		t.g.SetKeybinding("table", gocui.KeyEsc, gocui.ModNone, t.popTable),
		t.g.SetKeybinding("table", gocui.KeyArrowDown, gocui.ModNone, t.cursorDown),
		t.g.SetKeybinding("table", gocui.KeyArrowUp, gocui.ModNone, t.cursorUp),
		t.g.SetKeybinding("table", gocui.KeyEnter, gocui.ModNone, t.pushTableFromSelection),
		t.g.SetKeybinding("table", gocui.KeyPgup, gocui.ModNone, t.pageUp),
		t.g.SetKeybinding("table", gocui.KeyPgdn, gocui.ModNone, t.pageDown),
	)

	return err

}

func (t *TUI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (t *TUI) getSelectedY(v *gocui.View) int {
	_, y := v.Cursor()
	_, oy := v.Origin()

	return y + oy
}

func (t *TUI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	y := t.getSelectedY(v)
	if y < t.current.Length()-1 {
		v.MoveCursor(0, 1, false)
	}
	return nil
}

func (t *TUI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}

func (t *TUI) pageUp(g *gocui.Gui, v *gocui.View) error {

	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, vy := v.Size()
	if oy == 0 || oy+cy < vy {
		if err := v.SetOrigin(ox, 0); err != nil {
			return err
		}
	} else if oy <= vy {
		if err := v.SetOrigin(ox, oy+cy-vy); err != nil {
			return err
		}
	} else if err := v.SetOrigin(ox, oy-vy); err != nil {
		return err
	}
	if err := v.SetCursor(cx, 0); err != nil {
		return err
	}

	return nil
}

func (t *TUI) pageDown(g *gocui.Gui, v *gocui.View) error {

	ox, oy := v.Origin()
	cx, _ := v.Cursor()
	_, vy := v.Size()
	lr := len(v.BufferLines())
	if lr < vy {
		return nil
	}
	if oy+vy >= lr-vy {
		if err := v.SetOrigin(ox, lr-vy); err != nil {
			return err
		}
	} else if err := v.SetOrigin(ox, oy+vy); err != nil {
		return err
	}
	if err := v.SetCursor(cx, 0); err != nil {
		return err
	}

	return nil

}

func (t *TUI) pushTableFromSelection(g *gocui.Gui, v *gocui.View) error {
	y := t.getSelectedY(v)
	newTable, err := t.current.TableFromHere(y)
	newTable.SetGui(g)
	if newTable == nil || err != nil {
		return err
	}

	return t.pushTable(newTable, g, v)

}

func (t *TUI) pushTable(table ui.Table, g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	t.previousCursors = append(t.previousCursors, cy)
	t.previousOrigins = append(t.previousOrigins, oy)
	t.previousTables = append(t.previousTables, t.current)
	t.current = table

	err := v.SetCursor(0, 0)
	if err != nil {
		return err
	}
	err = v.SetOrigin(0, 0)
	if err != nil {
		return err
	}

	return nil
}

func (t *TUI) popTable(g *gocui.Gui, v *gocui.View) error {
	if len(t.previousTables) > 0 {
		lastCursor := t.previousCursors[len(t.previousCursors)-1]
		lastOrigin := t.previousOrigins[len(t.previousOrigins)-1]
		t.current = t.previousTables[len(t.previousTables)-1]

		t.previousCursors = t.previousCursors[:len(t.previousCursors)-1]
		t.previousOrigins = t.previousOrigins[:len(t.previousOrigins)-1]
		t.previousTables = t.previousTables[:len(t.previousTables)-1]

		err := v.SetCursor(0, lastCursor)
		if err != nil {
			return err
		}
		err = v.SetOrigin(0, lastOrigin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TUI) showClusterInfo(g *gocui.Gui, v *gocui.View) error {

	cluster := ui.NewClusterTable(t.cc.Admin)
	if t.current.Name() == cluster.Name() {
		return nil
	}

	return t.pushTable(cluster, g, v)
}

func (t *TUI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("header", -1, -1, maxX, 3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}

	v.Clear()

	t.current.Header(v, maxX)

	v, err = g.SetView("table", -1, 1, maxX, maxY-2)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorWhite
		v.SelFgColor = gocui.ColorBlack

		_, err = g.SetCurrentView("table")

		if err != nil {
			return err
		}
	}

	v.Clear()
	t.current.Render(g, v, maxX)

	v, err = g.SetView("footer", -1, maxY-2, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		v.BgColor = gocui.ColorBlack

	}

	serverAddress := fmt.Sprintf("Server [%s]", t.config.Address)
	keys := "[q] Quit [esc] Go back [s] Describe cluster"
	v.Clear()
	fmt.Fprintf(v, "%s%s", keys, fmt.Sprintf("%*s", maxX-len(keys), serverAddress))

	return nil
}
