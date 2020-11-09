package azurerm

import (
	"fmt"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/provider"

	tfdocAzurerm "github.com/cycloidio/tfdocs/providers/azurerm"
	"github.com/cycloidio/tfdocs/resource"
)

// Provider is the Type that defines the interface 'provider.Interface'
type Provider struct {
	provider.NopProvider
}

// Type returns the type of the implementation
func (a Provider) Type() provider.Type { return provider.Azurerm }

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
	r, err := tfdocAzurerm.GetResource(resource)
	if err != nil {
		return nil, fmt.Errorf("could not get resource %q: %w", resource, errcode.ErrProviderNotFoundResource)
	}
	return r, nil
}

// DataSource returns the resource information
func (a Provider) DataSource(resource string) (*resource.Resource, error) {
	r, err := tfdocAzurerm.GetDataSource(resource)
	if err != nil {
		return nil, fmt.Errorf("could not get resource %q: %w", resource, errcode.ErrProviderNotFoundDataSource)
	}
	return r, nil
}

// ResourceInOutNodes returns the In, Out and Nodes of the rs based on the cfg
func (a Provider) ResourceInOutNodes(id, rs string, cfgs map[string]map[string]interface{}) ([]string, []string, []string) {
	var ins, outs []string
	cfg := cfgs[id]
	switch rs {
	case "azurerm_virtual_network_peering":
		vnn := cfg["virtual_network_name"]
		vnID, ok := getRsIDByName(cfgs, vnn)
		if !ok {
			break
		}
		ins = append(ins, vnID)

		rvni := cfg["remote_virtual_network_id"]
		outs = append(outs, rvni.(string))
	}
	return ins, outs, nil
}

// getRsIDByName ranges over all resources, looking for a same name as provided.
// If found, the id of the resource will be returned.
func getRsIDByName(cfgs map[string]map[string]interface{}, name interface{}) (string, bool) {
	for _, cfg := range cfgs {
		rsName := cfg["name"]
		if rsName == name {
			pvid, ok := cfg["id"].(string)
			return pvid, ok
		}
	}
	return "", false
}

var (
	usedAttributes = []string{
		"remote_virtual_network_id",
		"virtual_network_name",
		"id",
		"name",
	}
)

// UsedAttributes returns all the attributes that are
// required/used/needed on the providers, so when we have to
// prune we know what to keep
func (a Provider) UsedAttributes() []string {
	return usedAttributes
}
