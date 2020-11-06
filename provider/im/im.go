package im

import (
	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/tfdocs/resource"
)

// Provider is the Type that defines the interface 'provider.Interface'
type Provider struct {
	provider.NopProvider
}

// Type returns the type of the implementation
func (a Provider) Type() provider.Type { return provider.IM }

// Resource returns the resource information
func (a Provider) Resource(rsc string) (*resource.Resource, error) {
	return &resource.Resource{
		Name: rsc,
		Type: rsc,
		Icon: "baseline_cloud_queue_black.svg",
	}, nil
}
