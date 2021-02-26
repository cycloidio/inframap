package aws

import (
	"fmt"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/provider"

	tfdocAWS "github.com/cycloidio/tfdocs/providers/aws"
	"github.com/cycloidio/tfdocs/resource"
)

// Provider is the Type that defines the interface 'provider.Interface'
type Provider struct {
	provider.NopProvider
}

var (
	usedAttributes = []string{
		"id",
		"egress",
		"ingress",
		"source_security_group_id",
		"security_group_id",
	}
)

// Type returns the type of the implementation
func (a Provider) Type() provider.Type { return provider.AWS }

// IsNode returns if the resources is a Node
func (a Provider) IsNode(resource string) bool {
	_, ok := nodeTypes[resource]
	return ok
}

// IsEdge returns if the resource is an Edge
func (a Provider) IsEdge(resource string) bool {
	_, ok := edgeTypes[resource]
	return ok
}

// Resource returns the resource information
func (a Provider) Resource(resource string) (*resource.Resource, error) {
	r, err := tfdocAWS.GetResource(resource)
	if err != nil {
		return nil, fmt.Errorf("could not get resource %q: %w", resource, errcode.ErrProviderNotFoundResource)
	}
	return r, nil
}

// DataSource returns the resource information
func (a Provider) DataSource(resource string) (*resource.Resource, error) {
	r, err := tfdocAWS.GetDataSource(resource)
	if err != nil {
		return nil, fmt.Errorf("could not get resource %q: %w", resource, errcode.ErrProviderNotFoundDataSource)
	}
	return r, nil
}

// ResourceInOutNodes returns the In, Out and Nodes of the rs based on the cfg
func (a Provider) ResourceInOutNodes(id, rs string, cfgs map[string]map[string]interface{}) ([]string, []string, []string) {
	var ins, outs, nodes []string
	cfg := cfgs[id]
	if rs == "aws_security_group" {
		ingress, inok := cfg["ingress"]
		if inok {
			for _, in := range ingress.([]interface{}) {
				min := in.(map[string]interface{})
				if sgs, ok := min["security_groups"].([]interface{}); ok {
					for _, sg := range sgs {
						ins = append(ins, sg.(string))
					}
				}
				if cidrs, ok := min["cidr_blocks"].([]interface{}); ok {
					for _, ci := range cidrs {
						sci := ci.(string)
						if sci == "0.0.0.0/0" {
							nodes = append(nodes, fmt.Sprintf("im_out.%v/%v->%v", min["protocol"], min["from_port"], min["to_port"]))
						}
					}
				}
			}
		}

		egress, egok := cfg["egress"]
		if egok {
			for _, eg := range egress.([]interface{}) {
				if sgs, ok := eg.(map[string]interface{})["security_groups"].([]interface{}); ok {
					for _, sg := range sgs {
						outs = append(outs, sg.(string))
					}
				}
			}
		}

	} else if rs == "aws_security_group_rule" {
		in, ok := cfg["source_security_group_id"]
		if ok && in != nil {
			ins = append(ins, in.(string))
		}

		out, ok := cfg["security_group_id"]
		if ok && in != nil {
			outs = append(outs, out.(string))
		}
	}

	return ins, outs, nodes
}

// UsedAttributes returns all the attributes that are
// required/used/needed on the providers, so when we have to
// prune we know what to keep
func (a Provider) UsedAttributes() []string {
	return usedAttributes
}
