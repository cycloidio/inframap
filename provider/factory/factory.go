package factory

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cycloidio/inframap/provider/azurerm"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/inframap/provider/aws"
	"github.com/cycloidio/inframap/provider/flexibleengine"
	"github.com/cycloidio/inframap/provider/google"
	"github.com/cycloidio/inframap/provider/im"
	"github.com/cycloidio/inframap/provider/openstack"
)

var (
	providers = map[provider.Type]provider.Provider{
		provider.IM:             im.Provider{},
		provider.AWS:            aws.Provider{},
		provider.FlexibleEngine: flexibleengine.Provider{},
		provider.OpenStack:      openstack.Provider{},
		provider.Google:         google.Provider{},
		provider.Azurerm:        azurerm.Provider{},
	}

	// reProvider is a regexp to match 'aws' from 'aws' or 'aws_iam_user'
	reProvider = regexp.MustCompile(`(?P<provider>[^_][a-z0-9]+)(?:_+)?`)
)

// GetProviderAndResource returns the Interface
// and the resource name "aws_alb.front" -> "aws_alb"
func GetProviderAndResource(can string) (provider.Provider, string, error) {
	// Due to modules, we'll check it from the back not from
	// the front as it may have modules prefix
	rss := strings.Split(can, ".")
	var rs string
	if len(rss) > 1 {
		rs = rss[len(rss)-2]
	} else {
		rs = can
	}
	match := reProvider.FindStringSubmatch(rs)

	if match == nil {
		return nil, "", fmt.Errorf("no match for %s: %w", can, errcode.ErrProviderNotFound)
	}

	t, err := provider.TypeString(match[1])
	if err != nil {
		return nil, "", fmt.Errorf("no provider type with %s from %s: %w", match[1], can, errcode.ErrProviderNotFound)
	}

	pv, ok := providers[t]
	if !ok {
		return nil, "", fmt.Errorf("no provider implementation for type %s from %s: %w", match[1], can, errcode.ErrProviderNotFound)
	}
	return pv, rs, nil
}
