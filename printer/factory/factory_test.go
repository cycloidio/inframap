package factory_test

import (
	"errors"
	"testing"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/printer/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		p, err := factory.Get("dot")
		require.NoError(t, err)
		assert.NotNil(t, p)
	})
	t.Run("Error", func(t *testing.T) {
		p, err := factory.Get("potato")
		assert.True(t, errors.Is(err, errcode.ErrPrinterNotFound))
		assert.Nil(t, p)
	})
}
