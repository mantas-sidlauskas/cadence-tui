package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"github.com/uber/cadence/client/admin"
)

var _ Table = (*Cluster)(nil)

type Cluster struct {
	Admin admin.Client
}

func (d *Cluster) Name() string {
	return "cluster"
}

func NewClusterTable(client admin.Client) Table {

	return &Cluster{
		Admin: client,
	}
}

var clusterCol = Columns{}.Add("info", -1)

func (d *Cluster) Header(v *gocui.View, maxX int) {
	color.New(color.FgWhite).Add(color.Bold).Fprintf(v, "Cluster info\n")
}

func (d *Cluster) Render(g *gocui.Gui, v *gocui.View, maxX int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := d.Admin.DescribeCluster(ctx)
	if err != nil {
		fmt.Fprintf(v, "Failed to describe cluster: %v\n", err)
		return
	}
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Fprintf(v, "Error when try to print pretty: %v\n", err)
		return
	}
	fmt.Fprintf(v, string(b))
}

func (d *Cluster) Length() int                                    { return 1000 }
func (d *Cluster) LoadNext() error                                { return nil }
func (d *Cluster) TableFromHere(selectedIndex int) (Table, error) { return nil, nil }
func (d *Cluster) SetGui(g *gocui.Gui)                            {}
