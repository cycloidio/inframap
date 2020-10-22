# Put here a custom name for the GKE Cluster
# Otherwise `${var.project}-${var.env}` will be used
locals {
  cluster_name = ""
}

# https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters
# You cannot use a cluster master, node, Pod, or Service IP range that overlaps with 172.17.0.0/16.
# The size of the RFC 1918 block for the cluster master must be /28.

module "vpc" {
  #####################################
  # Do not modify the following lines #
  source = "./module-vpc"

  project  = var.project
  env      = var.env
  customer = var.customer

  #####################################

  ###
  # General
  ###

  #. gcp_project (required):
  #+ The Google Cloud Platform project to use. 
  gcp_project = var.gcp_project

  #. gcp_region (optional): eu-central1
  #+ The Google Cloud Platform region to use.
  gcp_region = var.gcp_region

  #. extra_labels (optional): {}
  #+ Dict of extra labels to add on aws resources. format { "foo" = "bar" }.

  ###
  # Networking
  ###

  #. subnet_cidr (optional): 10.8.0.0/16
  #+ The CIDR of the VPC subnet.
  subnet_cidr = "10.8.0.0/16"

  #. pods_cidr (optional): 10.9.0.0/16
  #+ The CIDR of the pods secondary range.
  pods_cidr = "10.9.0.0/16"

  #. services_cidr (optional): 10.10.0.0/16
  #+ The CIDR of the services secondary range.
  services_cidr = "10.10.0.0/16"

  #. network_routing_mode (optional): GLOBAL
  #+ The network routing mode.

  ###
  # Required (should probably not be touched)
  ###

  cluster_name = local.gke_cluster_name
}

module "gke" {
  #####################################
  # Do not modify the following lines #
  source = "./module-gke"

  project  = var.project
  env      = var.env
  customer = var.customer

  #####################################

  ###
  # General
  ###

  #. gcp_project (required):
  #+ The Google Cloud Platform project to use. 
  gcp_project = var.gcp_project

  #. gcp_region (optional): eu-central1
  #+ The Google Cloud Platform region to use.
  gcp_region = var.gcp_region

  #. gcp_zones (optional): []
  #+ To use specific Google Cloud Platform zones if not regional, otherwise it will be chosen randomly.

  #. extra_labels (optional): {}
  #+ Dict of extra labels to add on GCP resources. format { "foo" = "bar" }.

  ###
  # Control plane
  ###

  #. cluster_version (optional): latest
  #+ GKE Cluster version to use.

  #. cluster_release_channel (optional): UNSPECIFIED
  #+ GKE Cluster release channel to use. Accepted values are UNSPECIFIED, RAPID, REGULAR and STABLE.

  #. cluster_regional (optional): false
  #+ If the GKE Cluster must be regional or zonal. Be careful, this setting is destructive.

  #. enable_only_private_endpoint (optional): false
  #+ If true, only enable the private endpoint which disable the Public endpoint entirely. If false, private endpoint will be enabled, and the public endpoint will be only accessible by master authorized networks.

  #. master_authorized_networks (optional): []
  #+ List of master authorized networks.
  # master_authorized_networks = [
  #   {
  #     name: "my-ip",
  #     cidr: "x.x.x.x/32"
  #   }
  # ]

  #. enable_network_policy (optional): true
  #+ Enable GKE Cluster network policies addon.

  #. enable_horizontal_pod_autoscaling (optional): true
  #+ Enable GKE Cluster horizontal pod autoscaling addon.

  #. enable_vertical_pod_autoscaling (optional): false
  #+ Enable GKE Cluster vertical pod autoscaling addon. Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it.

  #. enable_http_load_balancing (optional): false
  #+ Enable GKE Cluster HTTP load balancing addon.

  #. enable_binary_authorization (optional): false
  #+ Enable GKE Cluster BinAuthZ Admission controller.

  #. enable_cloudrun (optional): false
  #+ Enable GKE Cluster Cloud Run for Anthos addon.

  #. enable_istio (optional): false
  #+ Enable GKE Cluster Istio addon.

  ###
  # Node pools
  ###

  #. node_pools (optional): []
  #+ GKE Cluster node pools to create.
  node_pools = [
    {
      name         = "my-node-pool"
      machine_type = "n1-standard-1"
      image_type   = "COS"

      auto_repair  = true
      auto_upgrade = false
      preemptible  = false

      autoscaling        = true
      initial_node_count = 1
      min_count          = 1
      max_count          = 1

      # autoscaling = false
      # node_count = 1

      local_ssd_count = 0
      disk_size_gb    = 100
      disk_type       = "pd-ssd"

      # service_account = ""
      # accelerator_count = 0
      # accelerator_type = ""

      # oauth_scopes = []
      # metadata     = {}
      # labels       = {}
      # taints       = []
      # tags         = []
    },
  ]

  #. enable_shielded_nodes (optional): true
  #+ Enable GKE Cluster Shielded Nodes features on all nodes.

  #. enable_sandbox (optional): false
  #+ Enable GKE Sandbox (Do not forget to set image_type = COS_CONTAINERD and node_version = 1.12.7-gke.17 or later to use it).

  #. default_max_pods_per_node (optional): 110
  #+ The maximum number of pods to schedule per node.

  ###
  # Required (should probably not be touched)
  ###

  cluster_name      = local.gke_cluster_name
  subnet_name       = module.vpc.subnet_name
  pods_ip_range     = module.vpc.pods_ip_range
  services_ip_range = module.vpc.services_ip_range
}
