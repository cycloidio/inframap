package google_test

import (
	"fmt"
	"testing"

	"github.com/cycloidio/inframap/provider"
	"github.com/cycloidio/inframap/provider/google"
	"github.com/stretchr/testify/assert"
)

func TestResourceInOutNodes(t *testing.T) {
	t.Run("SuccessGCF_INGRESS", func(t *testing.T) {
		google := google.Provider{}
		id := "id"
		inid := "id-in"
		outid := "out-id"
		tagIn := "tag-in"
		tagOut := "tag-out"
		rs := "google_compute_firewall"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"direction": "INGRESS",
				"target_tags": []interface{}{
					tagOut,
				},
				"source_tags": []interface{}{
					tagIn,
				},
			},
			"in-node-id": map[string]interface{}{
				"id": inid,
				"tags": []interface{}{
					tagIn,
				},
			},
			"out-node-id": map[string]interface{}{
				"id": outid,
				"tags": []interface{}{
					tagOut,
				},
			},
		}

		ins, outs, _ := google.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{inid}, ins)
		assert.Equal(t, []string{outid}, outs)
	})
	t.Run("SuccessGCF_INGRESS_with_HCLCanonicalKey", func(t *testing.T) {
		google := google.Provider{}
		id := "id"
		incan := "id-can"
		outcan := "out-can"
		tagIn := "tag-in"
		tagOut := "tag-out"
		rs := "google_compute_firewall"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"direction": "INGRESS",
				"target_tags": []interface{}{
					tagOut,
				},
				"source_tags": []interface{}{
					tagIn,
				},
			},
			"in-node-id": map[string]interface{}{
				provider.HCLCanonicalKey: incan,
				"tags": []interface{}{
					tagIn,
				},
			},
			"out-node-id": map[string]interface{}{
				provider.HCLCanonicalKey: outcan,
				"tags": []interface{}{
					tagOut,
				},
			},
		}

		ins, outs, _ := google.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{fmt.Sprintf("${%s.something}", incan)}, ins)
		assert.Equal(t, []string{fmt.Sprintf("${%s.something}", outcan)}, outs)
	})
	t.Run("SuccessGCF_EGRESS", func(t *testing.T) {
		google := google.Provider{}
		id := "id"
		inid := "id-in"
		tagIn := "tag-in"
		rs := "google_compute_firewall"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"direction": "EGRESS",
				"target_tags": []interface{}{
					tagIn,
				},
			},
			"in-node-id": map[string]interface{}{
				"id": inid,
				"tags": []interface{}{
					tagIn,
				},
			},
		}

		ins, outs, _ := google.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{inid}, ins)
		assert.Equal(t, []string(nil), outs)
	})
	t.Run("SuccessGCF_EGRESS_with_HCLCanonicalKey", func(t *testing.T) {
		google := google.Provider{}
		id := "id"
		incan := "id-can"
		tagIn := "tag-in"
		rs := "google_compute_firewall"
		cfg := map[string]map[string]interface{}{
			id: map[string]interface{}{
				"direction": "EGRESS",
				"target_tags": []interface{}{
					tagIn,
				},
			},
			"in-node-id": map[string]interface{}{
				provider.HCLCanonicalKey: incan,
				"tags": []interface{}{
					tagIn,
				},
			},
		}

		ins, outs, _ := google.ResourceInOutNodes(id, rs, cfg)
		assert.Equal(t, []string{fmt.Sprintf("${%s.something}", incan)}, ins)
		assert.Equal(t, []string(nil), outs)
	})
}
