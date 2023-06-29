package ui

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/mantas-sidlauskas/cadence-tui/cadence"
	"github.com/uber/cadence/common/types"
)

var _ Table = (*Domains)(nil)

type Domains struct {
	List   []*types.DescribeDomainResponse
	Client *cadence.Client
	g      *gocui.Gui
}

func (d *Domains) Name() string {
	return "domains"
}

func (d *Domains) SetGui(g *gocui.Gui) {
	d.g = g
}

func NewDomainsTable(client *cadence.Client) *Domains {

	return &Domains{
		List:   client.GetDomains(),
		Client: client,
	}
}

var domainColumns = Columns{}.
	Add("name", 32).
	Add("global", 14).
	Add("description", -1)

func (d *Domains) Header(v *gocui.View, maxX int) {

	hName := PaddedString("Name", domainColumns.GetSize("name", maxX), 2)
	hGlobal := PaddedString("Is global?", domainColumns.GetSize("global", maxX), 2)
	hDescription := PaddedString("Description", domainColumns.GetSize("description", maxX), 2)

	color.New(color.FgWhite).Add(color.Bold).Fprintf(v, "Domains\n")
	color.New(color.Bold).Fprintf(v, "%s %s %s\n", hName, hGlobal, hDescription)
}

func (d *Domains) Render(g *gocui.Gui, v *gocui.View, maxX int) {

	for _, domain := range d.List {
		name := PaddedString(domain.DomainInfo.Name, domainColumns.GetSize("name", maxX), 2)
		global := PaddedString(fmt.Sprintf("%t", domain.IsGlobalDomain), domainColumns.GetSize("global", maxX), 2)
		description := PaddedString(domain.GetDomainInfo().Description, domainColumns.GetSize("description", maxX), 2)

		fmt.Fprintf(v, "\n%s %s %s", name, global, description)
	}
}

func (d *Domains) Length() int {
	return len(d.List)
}

func (d *Domains) LoadNext() error {

	return nil
}

func (d *Domains) TableFromHere(selectedIndex int) (Table, error) {

	domain := d.List[selectedIndex]

	workflows, err := d.Client.GetWorkflows(domain.DomainInfo.Name)
	if err != nil {
		return nil, err
	}

	return NewWorkflowsTable(d.Client, domain.DomainInfo.Name, workflows), nil
}
