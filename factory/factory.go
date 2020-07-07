package factory

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cycloidio/infraview/aws"
	"github.com/cycloidio/infraview/errcode"
	"github.com/cycloidio/infraview/provider"
)

var (
	providers = map[provider.Type]provider.Provider{
		provider.AWS: aws.Provider{},
	}

	// reProvider is a regexp to match 'aws' from 'aws' or 'aws_iam_user'
	reProvider = regexp.MustCompile(`(?P<provider>[^_][a-z0-9]+)(?:_+)?`)
)

// GetProviderAndResource returns the Interface
// and the resource name "aws_alb.front" -> "aws_alb"
func GetProviderAndResource(can string) (provider.Provider, string, error) {
	rs := strings.Split(can, ".")[0]
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
