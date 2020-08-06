package generate_test

import (
	"io/ioutil"
	"testing"

	"github.com/cycloidio/inframap/generate"
	"github.com/cycloidio/inframap/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromState_AWS(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/aws_state_sg.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "aws_lb.tQBgz",
				},
				&graph.Node{
					Canonical: "aws_launch_template.vIkyE",
				},
				&graph.Node{
					Canonical: "aws_db_instance.Cpbzf",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source:     "aws_lb.tQBgz",
					Target:     "aws_launch_template.vIkyE",
					Canonicals: []string{"aws_security_group.rZnGI", "aws_security_group.YPHPR"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.vIkyE",
					Target:     "aws_db_instance.Cpbzf",
					Canonicals: []string{"aws_security_group.YPHPR", "aws_security_group.LHwFh"},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})

	t.Run("SuccessSGR", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/aws_state_sgr.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "aws_elb.tMVdH",
				},
				&graph.Node{
					Canonical: "aws_instance.TObJL",
				},
				&graph.Node{
					Canonical: "aws_db_instance.qktIK",
				},
				&graph.Node{
					Canonical: "aws_elasticache_cluster.VUhMF",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source:     "aws_elb.tMVdH",
					Target:     "aws_instance.TObJL",
					Canonicals: []string{"aws_security_group.kuDkz", "aws_security_group_rule.pMOSN", "aws_security_group.UKblk"},
				},
				&graph.Edge{
					Source:     "aws_instance.TObJL",
					Target:     "aws_db_instance.qktIK",
					Canonicals: []string{"aws_security_group.mzSGd", "aws_security_group.kuDkz"},
				},
				&graph.Edge{
					Source:     "aws_instance.TObJL",
					Target:     "aws_elasticache_cluster.VUhMF",
					Canonicals: []string{"aws_security_group.KaWAd", "aws_security_group.kuDkz"},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})

	t.Run("MultipleHangingEdges", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/aws_state_multiple_hanging_edges.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "aws_lb.load-balancer",
				},
				&graph.Node{
					Canonical: "aws_ecs_service.wordpress",
				},
				&graph.Node{
					Canonical: "aws_ecs_cluster.ecs-cluster",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source: "aws_lb.load-balancer",
					Target: "aws_ecs_service.wordpress",
					Canonicals: []string{
						"aws_security_group.ecs_service",
						"aws_security_group.generated_aws_lb_load-balancer",
						"aws_security_group_rule.allow-alb",
					},
				},
				&graph.Edge{
					Source:     "aws_ecs_service.wordpress",
					Target:     "aws_ecs_cluster.ecs-cluster",
					Canonicals: []string(nil),
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})

	t.Run("WithCount", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/aws_state_with_count.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)
		assert.Len(t, g.Nodes, 2)
	})
	t.Run("SuccessVPC", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/aws_state_vpc.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "aws_elb.front",
					GroupIDs:  []string{"vpc-0d96ad69"},
				},
				&graph.Node{
					Canonical: "aws_instance.front",
					GroupIDs:  []string{"vpc-0d96ad69"},
				},
				&graph.Node{
					Canonical: "aws_db_instance.magento",
					GroupIDs:  []string{"vpc-0d96ad69"},
				},
				&graph.Node{
					Canonical: "aws_elasticache_cluster.redis",
					GroupIDs:  []string{"vpc-0d96ad69"},
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source: "aws_elb.front",
					Target: "aws_instance.front",
					Canonicals: []string{
						"aws_security_group.elb-front",
						"aws_security_group.front",
						"aws_security_group_rule.elb_to_front_http",
					},
				},
				&graph.Edge{
					Source: "aws_instance.front",
					Target: "aws_db_instance.magento",
					Canonicals: []string{
						"aws_security_group.front",
						"aws_security_group.rds",
					},
				},
				&graph.Edge{
					Source: "aws_instance.front",
					Target: "aws_elasticache_cluster.redis",
					Canonicals: []string{
						"aws_security_group.front",
						"aws_security_group.redis",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}

func TestFromState_OpenStack(t *testing.T) {
	t.Run("SuccessLB", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/openstack_state_lb.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "openstack_compute_instance_v2.AaCFA",
				},
				&graph.Node{
					Canonical: "openstack_compute_instance_v2.gZfYc",
				},
				&graph.Node{
					Canonical: "openstack_lb_loadbalancer_v2.PPdjL",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "openstack_compute_instance_v2.gZfYc",
					Source: "openstack_compute_instance_v2.AaCFA",
					Canonicals: []string{
						"openstack_networking_port_v2.GQStv",
						"openstack_networking_port_v2.PimKo",
						"openstack_networking_secgroup_rule_v2.uzQon",
						"openstack_networking_secgroup_v2.KFnza",
						"openstack_networking_secgroup_v2.ilWCI",
					},
				},
				&graph.Edge{
					Target: "openstack_compute_instance_v2.gZfYc",
					Source: "openstack_lb_loadbalancer_v2.PPdjL",
					Canonicals: []string{
						"openstack_lb_listener_v2.lzSKa",
						"openstack_lb_member_v2.FbSow",
						"openstack_lb_pool_v2.nwYyz",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
	t.Run("SuccessSG", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/openstack_state_sg.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "openstack_compute_instance_v2.sjSbA",
				},
				&graph.Node{
					Canonical: "openstack_compute_instance_v2.PiGtZ",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "openstack_compute_instance_v2.sjSbA",
					Source: "openstack_compute_instance_v2.PiGtZ",
					Canonicals: []string{
						"openstack_networking_port_v2.QWPcW",
						"openstack_networking_port_v2.tZVzk",
						"openstack_networking_secgroup_rule_v2.QdNte",
						"openstack_networking_secgroup_v2.DKQQX",
						"openstack_networking_secgroup_v2.ievUc",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})

}

func TestFromState_FlexibleEngine(t *testing.T) {
	t.Run("SuccessFlexibleEngine", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/flexibleengine_state.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.Aumwn",
				},
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.bkbLl",
				},
				&graph.Node{
					Canonical: "flexibleengine_blockstorage_volume_v2.hOHQu",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "flexibleengine_compute_instance_v2.Aumwn",
					Source: "flexibleengine_compute_instance_v2.bkbLl",
					Canonicals: []string{
						"flexibleengine_networking_port_v2.nLwOe",
						"flexibleengine_networking_port_v2.nLyuK",
						"flexibleengine_networking_secgroup_rule_v2.ZPPPO",
						"flexibleengine_networking_secgroup_v2.mDAul",
						"flexibleengine_networking_secgroup_v2.uiWdG",
					},
				},
				&graph.Edge{
					Target:     "flexibleengine_blockstorage_volume_v2.hOHQu",
					Source:     "flexibleengine_compute_instance_v2.bkbLl",
					Canonicals: nil,
				},
				&graph.Edge{
					Target:     "flexibleengine_blockstorage_volume_v2.hOHQu",
					Source:     "flexibleengine_compute_instance_v2.Aumwn",
					Canonicals: nil,
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})

	t.Run("SuccessWithFlexibleEngine011", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/flexibleengine_state_tf_011.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.CpCKR",
				},
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.uwGDt",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "flexibleengine_compute_instance_v2.CpCKR",
					Source: "flexibleengine_compute_instance_v2.uwGDt",
					Canonicals: []string{
						"flexibleengine_networking_port_v2.Maycs",
						"flexibleengine_networking_port_v2.QkWrX",
						"flexibleengine_networking_secgroup_rule_v2.CbxAB",
						"flexibleengine_networking_secgroup_v2.hlwMN",
						"flexibleengine_networking_secgroup_v2.wVneZ",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})

	t.Run("SuccessWithComputeInterfaceAttach", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/flexibleengine_state_attach.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.KnHdC",
				},
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.LJwaI",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "flexibleengine_compute_instance_v2.KnHdC",
					Source: "flexibleengine_compute_instance_v2.LJwaI",
					Canonicals: []string{
						"flexibleengine_compute_interface_attach_v2.cXqiJ",
						"flexibleengine_compute_interface_attach_v2.cvlHf",
						"flexibleengine_networking_port_v2.KMFHw",
						"flexibleengine_networking_port_v2.cXecO",
						"flexibleengine_networking_secgroup_rule_v2.mAaoN",
						"flexibleengine_networking_secgroup_v2.HjOXk",
						"flexibleengine_networking_secgroup_v2.tRwVT",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}

func TestFromState_Google(t *testing.T) {
	t.Run("SuccessGoogle", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/google_state.json")
		require.NoError(t, err)

		g, cfg, err := generate.FromState(src, generate.Options{Clean: true, Connections: true})
		require.NoError(t, err)
		require.NotNil(t, g)
		require.NotNil(t, cfg)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "google_compute_instance.ZthAT",
				},
				&graph.Node{
					Canonical: "google_compute_instance.lodiw",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "google_compute_instance.lodiw",
					Source: "google_compute_instance.ZthAT",
					Canonicals: []string{
						"google_compute_firewall.PEZjy",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}
