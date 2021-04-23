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

variable "network_routing_mode" {
  description = "The network routing mode."
  default     = "GLOBAL"
}

variable "subnet_cidr" {
  description = "The CIDR of the VPC subnet."
  default     = "10.8.0.0/16"
}

variable "pods_cidr" {
  description = "The CIDR of the pods secondary range."
  default     = "10.9.0.0/16"
}

variable "services_cidr" {
  description = "The CIDR of the services secondary range."
  default     = "10.10.0.0/16"
}

#
# Control plane
#

variable "cluster_name" {
  description = "EKS Cluster given name."
}