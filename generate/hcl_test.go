package generate_test

import (
	"testing"

	"github.com/cycloidio/inframap/generate"
	"github.com/cycloidio/inframap/graph"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFromHCL_AWS(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, err := generate.FromHCL(fs, "./testdata/aws_hcl_sg.tf", generate.Options{Clean: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "aws_lb.front",
				},
				&graph.Node{
					Canonical: "aws_launch_template.front",
				},
				&graph.Node{
					Canonical: "aws_db_instance.application",
				},
				&graph.Node{
					Canonical: "aws_elasticache_cluster.redis",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source:     "aws_lb.front",
					Target:     "aws_launch_template.front",
					Canonicals: []string{"aws_security_group.lb-front", "aws_security_group.front"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_db_instance.application",
					Canonicals: []string{"aws_security_group.front", "aws_security_group.rds"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_elasticache_cluster.redis",
					Canonicals: []string{"aws_security_group.redis", "aws_security_group.front"},
				},
			},
		}

		assertEqualGraph(t, eg, g, nil)
	})
}

func TestFromHCL_FlexibleEngine(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, err := generate.FromHCL(fs, "./testdata/flexibleengine_hcl.tf", generate.Options{Clean: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.instance_one",
				},
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.instance_two",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "flexibleengine_compute_instance_v2.instance_one",
					Source: "flexibleengine_compute_instance_v2.instance_two",
					Canonicals: []string{
						"flexibleengine_networking_port_v2.port_instance_two",
						"flexibleengine_networking_port_v2.port_instance_one",
						"flexibleengine_networking_secgroup_rule_v2.ssh_two_to_one",
						"flexibleengine_networking_secgroup_v2.secgroup_instance_two",
						"flexibleengine_networking_secgroup_v2.secgroup_instance_one",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, nil)
	})
}

func TestFromHCL_Module(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, err := generate.FromHCL(fs, "./testdata/tf-module/", generate.Options{Clean: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "aws_lb.front",
				},
				&graph.Node{
					Canonical: "aws_launch_template.front",
				},
				&graph.Node{
					Canonical: "aws_db_instance.application",
				},
				&graph.Node{
					Canonical: "aws_elasticache_cluster.redis",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source:     "aws_lb.front",
					Target:     "aws_launch_template.front",
					Canonicals: []string{"aws_security_group.lb-front", "aws_security_group.front"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_db_instance.application",
					Canonicals: []string{"aws_security_group.front", "aws_security_group.rds"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_elasticache_cluster.redis",
					Canonicals: []string{"aws_security_group.redis", "aws_security_group.front"},
				},
			},
		}

		assertEqualGraph(t, eg, g, nil)
	})
}
