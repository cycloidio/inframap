package provider

import (
	"github.com/cycloidio/tfdocs/resource"
)

// NopProvider holds the default methods for
// the provider.Interface so if one Provider
// does not implement one method we do not have
// to write the method
type NopProvider struct{}

// Type returns the name of the Provider
func (n NopProvider) Type() Type { return Type(9999) }

// IsNode checks if the resource should be considered
// a Node or not
func (n NopProvider) IsNode(rsc string) bool { return false }

// IsEdge checks if the resource should be considered
// an Edge or not
func (n NopProvider) IsEdge(rsc string) bool { return false }

// Resource returns the resource information
func (n NopProvider) Resource(rsc string) (*resource.Resource, error) { return nil, nil }

// DataSource returns the resource information
func (n NopProvider) DataSource(rsc string) (*resource.Resource, error) { return nil, nil }

// ResourceInOut returns the resource In Out from a
// state config. As an example in AWS this would be
// an "aws_security_group" "ingress" and "egress"
func (n NopProvider) ResourceInOut(id, rs string, cfgs map[string]map[string]interface{}) (in, out []string) {
	return nil, nil
}

// UsedAttributes returns all the attributes that are
// required/used/needed on the providers, so when we have to
// prune we know what to keep
func (n NopProvider) UsedAttributes() []string { return nil }

// PreProcess defines new edges from the config.
// each element is an edge and for each edge we have the source
// and the target.
// [_][0] is the source of the edge
// [_][1] is the target of the edge
func (n NopProvider) PreProcess(cfg map[string]map[string]interface{}) [][]string {
	return nil
}
