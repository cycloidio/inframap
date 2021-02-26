package aws

var (
	// nodeTypes all the Nodes that we support right now
	nodeTypes = map[string]struct{}{
		"aws_api_gateway_resource":          struct{}{},
		"aws_athena_database":               struct{}{},
		"aws_autoscaling_group":             struct{}{},
		"aws_dax_cluster":                   struct{}{},
		"aws_db_instance":                   struct{}{},
		"aws_directory_service_directory":   struct{}{},
		"aws_dms_replication_instance":      struct{}{},
		"aws_dx_gateway":                    struct{}{},
		"aws_dynamodb_table":                struct{}{},
		"aws_ebs_volume":                    struct{}{},
		"aws_ecs_cluster":                   struct{}{},
		"aws_ecs_service":                   struct{}{},
		"aws_efs_file_system":               struct{}{},
		"aws_eip":                           struct{}{},
		"aws_eks_cluster":                   struct{}{},
		"aws_elasticache_cluster":           struct{}{},
		"aws_elastic_beanstalk_application": struct{}{},
		"aws_elasticsearch_domain":          struct{}{},
		"aws_elb":                           struct{}{},
		"aws_emr_cluster":                   struct{}{},
		"aws_iam_user":                      struct{}{},
		"aws_instance":                      struct{}{},
		"aws_internet_gateway":              struct{}{},
		"aws_kinesis_stream":                struct{}{},
		"aws_lambda_function":               struct{}{},
		"aws_lightsail_instance":            struct{}{},
		"aws_mq_broker":                     struct{}{},
		"aws_media_store_container":         struct{}{},
		"aws_nat_gateway":                   struct{}{},
		"aws_rds_cluster":                   struct{}{},
		"aws_rds_cluster_instance":          struct{}{},
		"aws_redshift_cluster":              struct{}{},
		"aws_s3_bucket":                     struct{}{},
		"aws_storagegateway_gateway":        struct{}{},
		"aws_sqs_queue":                     struct{}{},
		"aws_vpn_gateway":                   struct{}{},
		"aws_batch_job_definition":          struct{}{},
		"aws_neptune_cluster":               struct{}{},
		"aws_alb":                           struct{}{},
		"aws_lb":                            struct{}{},
		"aws_cloudfront_distribution":       struct{}{},
		"aws_elasticache_replication_group": struct{}{},
		"aws_launch_template":               struct{}{},
	}

	// noSecurityGroup is a map of all the resourcese that do not
	// have/relay on SecurityGroups
	noSecurityGroup = map[string]struct{}{
		"aws_s3_bucket":               struct{}{},
		"aws_cloudfront_distribution": struct{}{},
	}

	// edgeTypes map of all the supported Edges
	edgeTypes = map[string]struct{}{
		"aws_security_group":      struct{}{},
		"aws_security_group_rule": struct{}{},
	}
)
