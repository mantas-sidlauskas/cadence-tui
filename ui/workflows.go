package ui

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/mantas-sidlauskas/cadence-tui/cadence"
	"github.com/uber/cadence/common/types"
)

var _ Table = (*Workflows)(nil)

type Workflows struct {
	workflows []*types.WorkflowExecutionInfo
	cadence   *cadence.Client
	domain    string
	g         *gocui.Gui
}

func (a *Workflows) Name() string {
	return "workflows"
}

func (a *Workflows) SetGui(g *gocui.Gui) {
	a.g = g
}

func NewWorkflowsTable(client *cadence.Client, domain string, workflows []*types.WorkflowExecutionInfo) *Workflows {
	return &Workflows{
		workflows: workflows,
		cadence:   client,
		domain:    domain,
	}
}

var workflowColumns = Columns{}.
	Add("wt", 22).
	Add("wid", 40).
	Add("rid", 40).
	Add("tl", 24).
	Add("cron", 10).
	Add("st", 34).
	Add("ext", 34).
	Add("endtime", 34)

func (a *Workflows) Header(v *gocui.View, maxX int) {

	wtH := PaddedString("WORKFLOW TYPE", workflowColumns.GetSize("wt", maxX), 2)
	widH := PaddedString("WORKFLOW ID", workflowColumns.GetSize("wid", maxX), 2)
	runH := PaddedString("RUN ID", workflowColumns.GetSize("rid", maxX), 2)
	tlH := PaddedString("TASK LIST", workflowColumns.GetSize("tl", maxX), 2)
	crH := PaddedString("Is Cron", workflowColumns.GetSize("cron", maxX), 2)
	stTH := PaddedString("Start time", workflowColumns.GetSize("st", maxX), 2)
	etT := PaddedString("Execution time", workflowColumns.GetSize("ext", maxX), 2)
	endT := PaddedString("End time", workflowColumns.GetSize("endtime", maxX), 2)

	color.New(color.FgWhite).Add(color.Bold).Fprintf(v, fmt.Sprintf("Workflows in domain %q\n", a.domain))
	color.New(color.Bold).Fprintf(v, "%s %s %s %s %s %s %s %s\n", wtH, widH, runH, tlH, crH, stTH, etT, endT)

}

func (a *Workflows) Render(g *gocui.Gui, v *gocui.View, maxX int) {

	for _, wf := range a.workflows {
		wt := PaddedString(wf.Type.Name, workflowColumns.GetSize("wt", maxX), 2)
		wid := PaddedString(wf.Execution.WorkflowID, workflowColumns.GetSize("wid", maxX), 2)
		rid := PaddedString(wf.Execution.RunID, workflowColumns.GetSize("rid", maxX), 2)
		tl := PaddedString(wf.TaskList, workflowColumns.GetSize("tl", maxX), 2)
		cron := PaddedString(fmt.Sprintf("%t", wf.IsCron), workflowColumns.GetSize("cron", maxX), 2)
		stTH := PaddedString(time.Unix(0, wf.GetStartTime()).Format(time.RFC3339), workflowColumns.GetSize("st", maxX), 2)
		etT := PaddedString(time.Unix(0, wf.GetExecutionTime()).Format(time.RFC3339), workflowColumns.GetSize("ext", maxX), 2)
		endT := PaddedString(time.Unix(0, wf.GetCloseTime()).Format(time.RFC3339), workflowColumns.GetSize("endtime", maxX), 2)

		fmt.Fprintf(v, "\n%s %s %s %s %s %s %s %s", wt, wid, rid, tl, cron, stTH, etT, endT)
	}
}

func (a *Workflows) Length() int {
	return len(a.workflows)
}

func (a *Workflows) LoadNext() error {
	return nil
}

func (a *Workflows) TableFromHere(selectedIndex int) (Table, error) {
	if len(a.workflows) < 1 {
		return nil, nil
	}
	execution := a.workflows[selectedIndex]
	return NewRunTable(a.g, a.cadence, a.domain, execution), nil
}
