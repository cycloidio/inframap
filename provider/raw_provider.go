package provider

import "github.com/cycloidio/tfdocs/resource"

type RawProvider struct {
	NopProvider
}

// Type returns the name of the Provider
func (n RawProvider) Type() Type { return Raw }

// IsNode checks if the resource should be considered
// a Node or not
// For the default one we consider all of the Nodes
func (n RawProvider) IsNode(rsc string) bool { return true }

// Resource returns the resource information
func (n RawProvider) Resource(rsc string) (*resource.Resource, error) {
	return &resource.Resource{
		Name: rsc,
		Type: rsc,
	}, nil
}
