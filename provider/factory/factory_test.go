package factory_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/inframap/provider/factory"
)

func TestGetProviderAndResource(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		eProvider provider.Type
		eResource string
		eError    error
	}{
		{
			name:      "Success",
			input:     "aws_lb.front",
			eProvider: provider.AWS,
			eResource: "aws_lb",
		},
		{
			name:      "SuccessOnlyResource",
			input:     "aws_lb",
			eProvider: provider.AWS,
			eResource: "aws_lb",
		},
		{
			name:   "ErrProviderNotFound-InvalidProvider",
			input:  "pepe_a.front",
			eError: errcode.ErrProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, rs, err := factory.GetProviderAndResource(tt.input)

			if tt.eError == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.eResource, rs)
				assert.Equal(t, tt.eProvider, p.Type())
			} else {
				assert.True(t, errors.Is(err, tt.eError))
			}
		})
	}
}
