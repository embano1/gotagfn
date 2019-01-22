package faastagger

import (
	"context"
	"log"
	"net/url"

	"github.com/vmware/govmomi/vapi/tags"

	"github.com/vmware/govmomi/vapi/rest"

	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi"
)

// Client is the tagging client used to tag a VM
type Client struct {
	vc *govmomi.Client
	// rest *rest.Client
	manager *tags.Manager
}

// New creates a tagging client
func New(ctx context.Context, vc string, insecure bool) (*Client, error) {
	u, err := url.Parse(vc)
	if err != nil {
		log.Printf("could not parse vCenter client URL: %v", err)
		return nil, err
	}

	c, err := govmomi.NewClient(ctx, u, insecure)
	if err != nil {
		log.Printf("could not get vCenter client: %v", err)
		return nil, err
	}

	r := rest.NewClient(c.Client)
	tm := tags.NewManager(r)

	return &Client{vc: c, manager: tm}, nil
}

// TagVM adds a tag to a virtual machine (ManagedObjectReference)
func (c *Client) TagVM(ctx context.Context, moref *types.ManagedObjectReference, tagID string) error {
	return c.manager.AttachTag(ctx, tagID, moref)
}
