package flexibleengine_test

import (
	"testing"

	"github.com/cycloidio/inframap/provider/flexibleengine"
	"github.com/stretchr/testify/assert"
)

func TestResourceInOutNodes(t *testing.T) {
	t.Run("SuccessCIAV2", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_compute_interface_attach_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"instance_id": "in-id",
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string(nil), outs)
	})
	t.Run("SuccessNSGRV2_ingress", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_networking_secgroup_rule_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"direction":       "ingress",
				"remote_group_id": "in-id",
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string(nil), outs)
	})
	t.Run("SuccessNSGRV2_egress", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_networking_secgroup_rule_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"direction":       "egress",
				"remote_group_id": "out-id",
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string(nil), ins)
		assert.Equal(t, []string{"out-id"}, outs)
	})
	t.Run("SuccessNPV2", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_networking_port_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"security_group_ids": []interface{}{
					"in-id",
				},
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string(nil), outs)
	})
	t.Run("SuccessLBLV2", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_lb_listener_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"loadbalancer_id": "in-id",
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string(nil), outs)
	})
	t.Run("SuccessLBPV2", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_lb_pool_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"listener_id": "in-id",
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string(nil), outs)
	})
	t.Run("SuccessLBMV2", func(t *testing.T) {
		fe := flexibleengine.Provider{}
		id := "id"
		rs := "flexibleengine_lb_member_v2"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"pool_id": "in-id",
			},
		}

		ins, outs, _ := fe.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string(nil), outs)
	})
}
