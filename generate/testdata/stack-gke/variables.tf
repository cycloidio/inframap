# Cycloid requirements
variable "project" {
  description = "Cycloid project name."
}

variable "env" {
  description = "Cycloid environment name."
}

variable "customer" {
  description = "Cycloid customer name."
}

# GCP
variable "gcp_project" {
  description = "GCP project to launch services."
  default     = "kubernetes-gke"
}

variable "gcp_region" {
  description = "GCP region to launch services."
  default     = "europe-west1"
}

# EKS
locals {
  gke_cluster_name = length(local.cluster_name) > 0 ? local.cluster_name : "${var.project}-${var.env}"
}