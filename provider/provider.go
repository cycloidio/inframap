package provider

import (
	"github.com/cycloidio/tfdocs/resource"
)

// Provider is an interface to abstract common functions on all the
// providers
type Provider interface {
	// Type returns the name of the Provider
	Type() Type

	// IsNode checks if the resource should be considered
	// a Node or not
	IsNode(rsc string) bool

	// IsEdge checks if the resource should be considered
	// an Edge or not
	IsEdge(rsc string) bool

	// Resource returns the resource information
	Resource(rsc string) (*resource.Resource, error)

	// DataSource returns the resource information
	DataSource(rsc string) (*resource.Resource, error)

	// ResourceInOut returns the resource In Out from a
	// state config. As an example in AWS this would be
	// an "aws_security_group" "ingress" and "egress"
	ResourceInOut(id, rs string, cfg map[string]map[string]interface{}) (in, out []string)

	// UsedAttributes returns all the attributes that are
	// required/used/needed on the providers, so when we have to
	// prune we know what to keep
	UsedAttributes() []string

	// PreProcess defines new edges from the config.
	// each element is an edge and for each edge we have the source
	// and the target.
	// [_][0] is the source of the edge
	// [_][1] is the target of the edge
	PreProcess(cfg map[string]map[string]interface{}) [][]string
}
