package flexibleengine

import (
	"fmt"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/provider"

	tfdocFE "github.com/cycloidio/tfdocs/providers/flexibleengine"
	"github.com/cycloidio/tfdocs/resource"
)

// Provider is the Type that defines the interface 'provider.Interface'
type Provider struct {
	provider.NopProvider
}

var (
	usedAttributes = []string{
		"id",
		"instance_id",
		"direction",
		"remote_group_id",
		"security_group_ids",
		"loadbalancer_id",
		"listener_id",
		"pool_id",
	}
)

// Type returns the type of the implementation
func (a Provider) Type() provider.Type { return provider.FlexibleEngine }

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
	r, err := tfdocFE.GetResource(resource)
	if err != nil {
		return nil, fmt.Errorf("could not get resource %q: %w", resource, errcode.ErrProviderNotFoundResource)
	}
	return r, nil
}

// DataSource returns the resource information
func (a Provider) DataSource(resource string) (*resource.Resource, error) {
	r, err := tfdocFE.GetDataSource(resource)
	if err != nil {
		return nil, fmt.Errorf("could not get resource %q: %w", resource, errcode.ErrProviderNotFoundDataSource)
	}
	return r, nil
}

// ResourceInOut returns the In and Out of the rs based on the cfg
func (a Provider) ResourceInOut(id, rs string, cfgs map[string]map[string]interface{}) ([]string, []string) {
	var ins, outs []string
	cfg := cfgs[id]
	switch rs {
	case "flexibleengine_compute_interface_attach_v2":
		if instanceID, ok := cfg["instance_id"]; ok {
			ins = append(ins, instanceID.(string))
		}
	case "flexibleengine_networking_secgroup_rule_v2":
		if direction, ok := cfg["direction"]; ok {
			if direction.(string) == "ingress" {
				if rg, ok := cfg["remote_group_id"]; ok {
					ins = append(ins, rg.(string))
				}
			} else {
				if rg, ok := cfg["remote_group_id"]; ok {
					outs = append(outs, rg.(string))
				}
			}
		}
	case "flexibleengine_networking_port_v2":
		if sgs, ok := cfg["security_group_ids"]; ok {
			for _, sg := range sgs.([]interface{}) {
				ins = append(ins, sg.(string))
			}
		}
	case "flexibleengine_lb_listener_v2":
		if lbid, ok := cfg["loadbalancer_id"]; ok {
			ins = append(ins, lbid.(string))
		}
	case "flexibleengine_lb_pool_v2":
		if lid, ok := cfg["listener_id"]; ok {
			ins = append(ins, lid.(string))
		}
	case "flexibleengine_lb_member_v2":
		if pid, ok := cfg["pool_id"]; ok {
			ins = append(ins, pid.(string))
		}
	}
	return ins, outs
}

// UsedAttributes returns all the attributes that are
// required/used/needed on the providers, so when we have to
// prune we know what to keep
func (a Provider) UsedAttributes() []string {
	return usedAttributes
}
