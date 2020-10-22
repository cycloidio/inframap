output "network_name" {
  description = "GKE Cluster dedicated network name."
  value       = module.gcp-network.network_name
}

output "network_self_link" {
  description = "GKE Cluster dedicated network URI."
  value       = module.gcp-network.network_self_link
}

output "subnet_name" {
  description = "GKE Cluster dedicated subnet name."
  value       = module.gcp-network.subnets_names[0]
}

output "subnet_self_link" {
  description = "GKE Cluster dedicated subnet URI."
  value       = module.gcp-network.subnets_self_links[0]
}

output "subnet_region" {
  description = "GKE Cluster dedicated subnet region."
  value       = module.gcp-network.subnets_regions[0]
}

output "pods_ip_range" {
  description = "GKE Cluster dedicated pods IP range."
  value       = module.gcp-network.subnets_secondary_ranges[0].*.range_name[0]
}

output "services_ip_range" {
  description = "GKE Cluster dedicated services IP range."
  value       = module.gcp-network.subnets_secondary_ranges[0].*.range_name[1]
}
