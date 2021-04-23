#
# General
#

variable "gcp_project" {
  description = "The Google Cloud Platform project to use."
}

variable "gcp_region" {
  description = "The Google Cloud Platform region to use."
  default     = "eu-central1"
}

variable "gcp_zones" {
  description = "To use specific Google Cloud Platform zones if not regional, otherwise it will be chosen randomly."
  default     = []
}

data "google_compute_zones" "available" {
  project = var.gcp_project
  region  = var.gcp_region
  status  = "UP"
}

locals {
  gcp_available_zones = length(var.gcp_zones) > 0 ? var.gcp_zones : data.google_compute_zones.available.names
}

variable "project" {
  description = "Cycloid project name."
}

variable "env" {
  description = "Cycloid environment name."
}

variable "customer" {
  description = "Cycloid customer name."
}

variable "extra_labels" {
  description = "Extra labels to add to all resources."
  default     = {}
}

locals {
  standard_labels = {
    cycloidio    = "true"
    env          = var.env
    project      = var.project
    client       = var.customer
  }
  merged_labels = merge(local.standard_labels, var.extra_labels)
}

#
# Networking
#

variable "subnet_name" {
  description = "GKE Cluster subnet name to use."
}

variable "pods_ip_range" {
  description = "GKE Cluster pods IP range to use."
}

variable "services_ip_range" {
  description = "GKE Cluster services IP range to use."
}

variable "master_cidr" {
  description = "GKE Cluster masters IP CIDR to use."
  default     = "172.16.0.0/28"
}

#
# Control plane
#

variable "cluster_name" {
  description = "GKE Cluster given name."
}

variable "cluster_version" {
  description = "GKE Cluster version to use."
  default     = "latest"
}

variable "cluster_release_channel" {
  description = "GKE Cluster release channel to use. Accepted values are UNSPECIFIED, RAPID, REGULAR and STABLE."
  default     = "UNSPECIFIED"
}

variable "cluster_regional" {
  description = "If the GKE Cluster must be regional or zonal. Be careful, this setting is destructive."
  default     = false
}

variable "enable_only_private_endpoint" {
  description = "If true, only enable the private endpoint which disable the Public endpoint entirely. If false, private endpoint will be enabled, and the public endpoint will be only accessible by master authorized networks."
  default     = false
}

variable "grant_registry_access" {
  description = "Grants created cluster-specific service account storage.objectViewer role."
  default     = true
}

variable "master_authorized_networks" {
  description = "List of master authorized networks."
  default     = []
}

variable "enable_network_policy" {
  description = "Enable GKE Cluster network policies addon."
  default     = true
}

variable "network_policy_provider" {
  description = "The GKE Cluster network policies addon provider."
  default     = "CALICO"
}

variable "enable_horizontal_pod_autoscaling" {
  description = "Enable GKE Cluster horizontal pod autoscaling addon."
  default     = true
}

variable "enable_vertical_pod_autoscaling" {
  description = "Enable GKE Cluster vertical pod autoscaling addon. Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it."
  default     = false
}

variable "enable_http_load_balancing" {
  description = "Enable GKE Cluster HTTP load balancing addon."
  default     = false
}

variable "disable_legacy_metadata_endpoints" {
  description = "Disable GKE Cluster legacy metadata endpoints."
  default     = true
}

variable "enable_binary_authorization" {
  description = "Enable GKE Cluster BinAuthZ Admission controller."
  default     = false
}

variable "enable_cloudrun" {
  description = "Enable GKE Cluster Cloud Run for Anthos addon."
  default     = false
}

variable "enable_istio" {
  description = "Enable GKE Cluster Istio addon."
  default     = false
}

variable "maintenance_start_time" {
  description = "Time window specified for daily maintenance operations in RFC3339 format."
  default     = "05:00"
}

#
# Node pools
#

variable "node_pools" {
  description = "GKE Cluster node pools to create."
  default     = []
}

variable "enable_shielded_nodes" {
  description = "Enable GKE Cluster Shielded Nodes features on all nodes."
  default     = true
}

variable "enable_sandbox" {
  description = "Enable GKE Sandbox (Do not forget to set image_type = COS_CONTAINERD and node_version = 1.12.7-gke.17 or later to use it)."
  default     = false
}

variable "default_max_pods_per_node" {
  description = "The maximum number of pods to schedule per node."
  default     = "110"
}

variable "enable_intranode_visibility" {
  description = "Whether Intra-node visibility is enabled for this cluster. This makes same node pod to pod traffic visible for VPC network."
  default     = false
}
