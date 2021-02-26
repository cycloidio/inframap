package aws_test

import (
	"testing"

	"github.com/cycloidio/inframap/provider/aws"
	"github.com/stretchr/testify/assert"
)

func TestResourceInOutNodes(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		aws := aws.Provider{}
		id := "id"
		rs := "aws_security_group"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"ingress": []interface{}{
					map[string]interface{}{
						"security_groups": []interface{}{
							"in-id",
						},
					},
				},
				"egress": []interface{}{
					map[string]interface{}{
						"security_groups": []interface{}{
							"out-id",
						},
					},
				},
			},
		}

		ins, outs, _ := aws.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string{"out-id"}, outs)
	})
	t.Run("SuccessSGR", func(t *testing.T) {
		aws := aws.Provider{}
		id := "id"
		rs := "aws_security_group_rule"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"source_security_group_id": "in-id",
				"security_group_id":        "out-id",
			},
		}

		ins, outs, _ := aws.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{"in-id"}, ins)
		assert.Equal(t, []string{"out-id"}, outs)
	})
}
