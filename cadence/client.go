package cadence

import (
	"context"
	"time"

	adminv1 "github.com/uber/cadence-idl/go/proto/admin/v1"
	apiv1 "github.com/uber/cadence-idl/go/proto/api/v1"
	"github.com/uber/cadence/client/admin"
	"github.com/uber/cadence/client/frontend"
	"github.com/uber/cadence/common/types"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/transport/grpc"
)

type Client struct {
	d     *yarpc.Dispatcher
	fc    frontend.Client
	Admin admin.Client
}

func NewClient(address string) *Client {
	var outbounds transport.Outbounds
	outbounds = transport.Outbounds{Unary: grpc.NewTransport().NewSingleOutbound(address)}

	d := yarpc.NewDispatcher(yarpc.Config{
		Name:      "cadence-client",
		Outbounds: yarpc.Outbounds{"cadence-frontend": outbounds},
	})

	if err := d.Start(); err != nil {
		d.Stop()
	}
	cc := d.ClientConfig("cadence-frontend")
	fc := frontend.NewGRPCClient(
		apiv1.NewDomainAPIYARPCClient(cc),
		apiv1.NewWorkflowAPIYARPCClient(cc),
		apiv1.NewWorkerAPIYARPCClient(cc),
		apiv1.NewVisibilityAPIYARPCClient(cc),
	)

	return &Client{
		d:     d,
		fc:    fc,
		Admin: admin.NewGRPCClient(adminv1.NewAdminAPIYARPCClient(cc)),
	}
}

func (c *Client) GetDomains() []*types.DescribeDomainResponse {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var token []byte
	pagesize := int32(100)
	var res []*types.DescribeDomainResponse

	for more := true; more; more = len(token) > 0 {
		listRequest := &types.ListDomainsRequest{
			PageSize:      pagesize,
			NextPageToken: token,
		}
		cl, err := c.fc.ListDomains(ctx, listRequest)
		if err != nil {
			panic(err)
		}
		token = cl.GetNextPageToken()
		res = append(res, cl.GetDomains()...)
	}
	return res

}

func (c *Client) GetWorkflows(domain string) ([]*types.WorkflowExecutionInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.fc.ListWorkflowExecutions(ctx, &types.ListWorkflowExecutionsRequest{
		Domain: domain,
	})

	if err != nil {
		return nil, err
	}

	return r.GetExecutions(), nil
}

func (c *Client) FC() frontend.Client {
	return c.fc
}
