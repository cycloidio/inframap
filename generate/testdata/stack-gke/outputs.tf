# VPC
output "network_name" {
  description = "GKE Cluster dedicated network name."
  value       = module.vpc.network_name
}

output "network_self_link" {
  description = "GKE Cluster dedicated network URI."
  value       = module.vpc.network_self_link
}

output "subnet_name" {
  description = "GKE Cluster dedicated subnet name."
  value       = module.vpc.subnet_name
}

output "subnet_self_link" {
  description = "GKE Cluster dedicated subnet URI."
  value       = module.vpc.subnet_self_link
}

output "subnet_region" {
  description = "GKE Cluster dedicated subnet region."
  value       = module.vpc.subnet_region
}

output "pods_ip_range" {
  description = "GKE Cluster dedicated pods IP range."
  value       = module.vpc.pods_ip_range
}

output "services_ip_range" {
  description = "GKE Cluster dedicated services IP range."
  value       = module.vpc.services_ip_range
}


# EKS Cluster
output "cluster_name" {
  description = "GKE Cluster name."
  value       = module.gke.cluster_name
}

output "cluster_type" {
  description = "GKE Cluster type."
  value       = module.gke.cluster_type
}

output "cluster_location" {
  description = "GKE Cluster location (region if regional cluster, zone if zonal cluster)."
  value       = module.gke.cluster_location
}

output "cluster_region" {
  description = "GKE Cluster region."
  value       = module.gke.cluster_region
}

output "cluster_zones" {
  description = "GKE Cluster zones."
  value       = module.gke.cluster_zones
}

output "cluster_master_version" {
  description = "GKE Cluster master version."
  value       = module.gke.cluster_master_version
}

output "cluster_release_channel" {
  description = "GKE Cluster release channel."
  value       = module.gke.cluster_release_channel
}

output "cluster_master_authorized_networks_config" {
  description = "GKE Cluster networks from which access to master is permitted."
  value       = module.gke.cluster_master_authorized_networks_config
}

output "cluster_endpoint" {
  description = "GKE Cluster endpoint."
  value       = module.gke.cluster_endpoint
}

output "cluster_ca" {
  description = "GKE Cluster certificate authority."
  value       = module.gke.cluster_ca
}

output "node_pools_names" {
  description = "GKE Cluster node pools names."
  value       = module.gke.node_pools_names
}

output "node_pools_versions" {
  description = "GKE Cluster node pools versions."
  value       = module.gke.node_pools_versions
}

output "node_pools_service_account" {
  description = "GKE Cluster nodes default service account if not overriden in `node_pools`."
  value       = module.gke.node_pools_service_account
}

output "kubeconfig" {
  description = "Kubernetes config to connect to the GKE cluster."
  value       = module.gke.kubeconfig
}
