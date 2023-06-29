package ui

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/mantas-sidlauskas/cadence-tui/cadence"
	"github.com/uber/cadence/common/pagination"
	"github.com/uber/cadence/common/types"
	"github.com/uber/cadence/tools/cli"
	"time"
)

var _ Table = (*WorkflowRun)(nil)

type WorkflowRun struct {
	run      *types.WorkflowExecutionInfo
	client   *cadence.Client
	domain   string
	iterator pagination.Iterator
	entries  []string
}

func (w *WorkflowRun) Name() string {
	return "run"
}

func (w *WorkflowRun) SetGui(g *gocui.Gui) {
	// noop
}

var runColumns = Columns{}.
	Add("id", 4).
	Add("timestamp", 11).
	Add("event", 32).
	Add("details", -1)

func (w *WorkflowRun) Render(g *gocui.Gui, v *gocui.View, maxX int) {
	// show collected entries
	for _, l := range w.entries {
		fmt.Fprint(v, l)
	}
}

func (w *WorkflowRun) Header(v *gocui.View, maxX int) {

	idHeader := PaddedString("ID", runColumns.GetSize("id", maxX), 2)
	tsHeader := PaddedString("Timestamp", runColumns.GetSize("timestamp", maxX), 2)
	eventHeader := PaddedString("EVENT", runColumns.GetSize("event", maxX), 2)
	detailsHeader := PaddedString("DETAILS", runColumns.GetSize("details", maxX), 2)
	color.New(color.FgWhite).Add(color.Bold).Fprintf(v, "Domain %q WF ID %s Run: %s\n", w.domain, w.run.Execution.WorkflowID, w.run.Execution.RunID)
	color.New(color.Bold).Fprintf(v, "%s %s %s %s\n", idHeader, tsHeader, eventHeader, detailsHeader)
}

func (w *WorkflowRun) Length() int {
	return len(w.entries)

}

func (w *WorkflowRun) LoadNext() error {
	return nil
}

func (w *WorkflowRun) TableFromHere(selectedIndex int) (Table, error) {
	return nil, nil
}

func NewRunTable(g *gocui.Gui, client *cadence.Client, domain string, we *types.WorkflowExecutionInfo) Table {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := &WorkflowRun{
		run:    we,
		client: client,
		domain: domain,
	}
	iterator, err := cli.GetWorkflowHistoryIterator(ctx, client.FC(), domain, we.Execution.WorkflowID, we.Execution.RunID, true, nil)
	if err != nil {
		w.entries = append(w.entries, err.Error())
	}
	w.iterator = iterator
	maxX, _ := g.Size()
	go func() {
		for w.iterator.HasNext() {
			entity, _ := w.iterator.Next()
			e := entity.(*types.HistoryEvent)
			cF := cli.EventColorFunction(*e.EventType)
			event := PaddedString(e.EventType.String(), runColumns.GetSize("event", maxX), 2)
			id := PaddedString(fmt.Sprintf("%d", e.ID), runColumns.GetSize("id", maxX), 2)
			ts := PaddedString(time.Unix(0, *e.Timestamp).Format(time.TimeOnly), runColumns.GetSize("timestamp", maxX), 2)
			details := PaddedString(cli.HistoryEventToString(e, false, runColumns.GetSize("details", maxX)), runColumns.GetSize("details", maxX), 2)
			w.entries = append(w.entries, fmt.Sprintf("\n%s %s %s %s ", id, ts, cF(event), details))

			g.Update(func(g *gocui.Gui) error {
				return nil
			})
		}
	}()

	return w
}
