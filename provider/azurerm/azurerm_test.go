package azurerm_test

import (
	"testing"

	"github.com/cycloidio/inframap/provider/azurerm"
	"github.com/stretchr/testify/assert"
)

func TestResourceInOut(t *testing.T) {
	t.Run("SuccessVNetwork", func(t *testing.T) {
		aws := azurerm.Provider{}
		id := "id"
		rs := "azurerm_virtual_network_peering"
		cfg := map[string]map[string]interface{}{
			"virtual_network": {
				"name": "vnname",
				"id":   "src_v_network",
			},
			id: {
				"virtual_network_name":      "vnname",
				"remote_virtual_network_id": "remote_v_network",
			},
		}

		ins, outs := aws.ResourceInOut(id, rs, cfg)
		assert.Equal(t, []string{"src_v_network"}, ins)
		assert.Equal(t, []string{"remote_v_network"}, outs)
	})
}
