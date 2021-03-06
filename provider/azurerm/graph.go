package azurerm

var (
	nodeTypes = map[string]struct{}{
		"azurerm_app_service":                       struct{}{},
		"azurerm_app_service_certificate":           struct{}{},
		"azurerm_app_service_environment":           struct{}{},
		"azurerm_app_service_plan":                  struct{}{},
		"azurerm_application_gateway":               struct{}{},
		"azurerm_bastion_host":                      struct{}{},
		"azurerm_batch_account":                     struct{}{},
		"azurerm_batch_application":                 struct{}{},
		"azurerm_batch_pool":                        struct{}{},
		"azurerm_cdn_endpoint":                      struct{}{},
		"azurerm_cdn_profile":                       struct{}{},
		"azurerm_container_group":                   struct{}{},
		"azurerm_container_registry":                struct{}{},
		"azurerm_cosmosdb_account":                  struct{}{},
		"azurerm_data_factory":                      struct{}{},
		"azurerm_dedicated_host":                    struct{}{},
		"azurerm_dedicated_host_group":              struct{}{},
		"azurerm_dns_zone":                          struct{}{},
		"azurerm_firewall":                          struct{}{},
		"azurerm_frontdoor":                         struct{}{},
		"azurerm_function_app":                      struct{}{},
		"azurerm_image":                             struct{}{},
		"azurerm_kubernetes_cluster":                struct{}{},
		"azurerm_kubernetes_cluster_node_pool":      struct{}{},
		"azurerm_lb":                                struct{}{},
		"azurerm_linux_virtual_machine":             struct{}{},
		"azurerm_mariadb_database":                  struct{}{},
		"azurerm_mariadb_server":                    struct{}{},
		"azurerm_mssql_database":                    struct{}{},
		"azurerm_mssql_elasticpool":                 struct{}{},
		"azurerm_mssql_server":                      struct{}{},
		"azurerm_mssql_virtual_machine":             struct{}{},
		"azurerm_mysql_database":                    struct{}{},
		"azurerm_mysql_server":                      struct{}{},
		"azurerm_nat_gateway":                       struct{}{},
		"azurerm_nat_gateway_public_ip_association": struct{}{},
		"azurerm_netapp_account":                    struct{}{},
		"azurerm_netapp_pool":                       struct{}{},
		"azurerm_netapp_volume":                     struct{}{},
		"azurerm_postgresql_database":               struct{}{},
		"azurerm_postgresql_server":                 struct{}{},
		"azurerm_private_dns_zone":                  struct{}{},
		"azurerm_private_endpoint":                  struct{}{},
		"azurerm_public_ip":                         struct{}{},
		"azurerm_redis_cache":                       struct{}{},
		"azurerm_sql_database":                      struct{}{},
		"azurerm_sql_server":                        struct{}{},
		"azurerm_storage_container":                 struct{}{},
		"azurerm_virtual_machine":                   struct{}{},
		"azurerm_virtual_network_gateway":           struct{}{},
		"azurerm_vpn_gateway":                       struct{}{},
		"azurerm_windows_virtual_machine":           struct{}{},
		"azurerm_virtual_network":                   struct{}{},
	}

	// edgeTypes map of all the supported Edges
	edgeTypes = map[string]struct{}{
		"azurerm_virtual_network_peering": {},
	}
)
