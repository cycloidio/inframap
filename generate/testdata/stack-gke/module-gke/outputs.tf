output "cluster_name" {
  description = "GKE Cluster name."
  value       = module.gcp-gke.name
}

output "cluster_type" {
  description = "GKE Cluster type."
  value       = module.gcp-gke.type
}

output "cluster_location" {
  description = "GKE Cluster location (region if regional cluster, zone if zonal cluster)."
  value       = module.gcp-gke.location
}

output "cluster_region" {
  description = "GKE Cluster region."
  value       = module.gcp-gke.region
}

output "cluster_zones" {
  description = "GKE Cluster zones."
  value       = module.gcp-gke.zones
}

output "cluster_master_version" {
  description = "GKE Cluster master version."
  value       = module.gcp-gke.master_version
}

output "cluster_release_channel" {
  description = "GKE Cluster release channel."
  value       = module.gcp-gke.release_channel
}

output "cluster_endpoint" {
  description = "GKE Cluster endpoint."
  value       = "https://${module.gcp-gke.endpoint}"
}

output "cluster_ca" {
  description = "GKE Cluster certificate authority."
  value       = module.gcp-gke.ca_certificate
}

output "cluster_master_authorized_networks_config" {
  description = "GKE Cluster networks from which access to master is permitted."
  value       = module.gcp-gke.master_authorized_networks_config
}

output "node_pools_names" {
  description = "GKE Cluster node pools names."
  value       = module.gcp-gke.node_pools_names
}

output "node_pools_versions" {
  description = "GKE Cluster node pools versions."
  value       = module.gcp-gke.node_pools_versions
}

output "node_pools_service_account" {
  description = "GKE Cluster nodes default service account if not overriden in `node_pools`."
  value       = module.gcp-gke.service_account
}

locals {
  kubeconfig = <<KUBECONFIG


apiVersion: v1
clusters:
- cluster:
    server: https://${module.gcp-gke.endpoint}
    certificate-authority-data: ${module.gcp-gke.ca_certificate}
  name: gke-${var.cluster_name}
contexts:
- context:
    cluster: gke-${var.cluster_name}
    user: gcp-${var.cluster_name}
  name: ${var.cluster_name}
current-context: ${var.cluster_name}
kind: Config
preferences: {}
users:
- name: gcp-${var.cluster_name}
  user:    
    auth-provider:    
      config:    
        cmd-args: config config-helper --format=json    
        cmd-path: gcloud    
        expiry-key: '{.credential.token_expiry}'    
        token-key: '{.credential.access_token}'    
      name: gcp
KUBECONFIG
}

output "kubeconfig" {
  description = "Kubernetes config to connect to the GKE Cluster."
  value       = local.kubeconfig
}
