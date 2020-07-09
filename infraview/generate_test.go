package infraview_test

import (
	"io/ioutil"
	"testing"

	"github.com/cycloidio/infraview/graph"
	"github.com/cycloidio/infraview/infraview"
	"github.com/stretchr/testify/require"
)

func TestFromState_AWS(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		src, err := ioutil.ReadFile("./testdata/aws_sg.json")
		require.NoError(t, err)

		g, cfg, err := infraview.FromState(src)
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
		src, err := ioutil.ReadFile("./testdata/aws_sgr.json")
		require.NoError(t, err)

		g, cfg, err := infraview.FromState(src)
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
}

func TestFromState_OpenStack(t *testing.T) {
}

func TestFromState_FlexibelEngine(t *testing.T) {
}
