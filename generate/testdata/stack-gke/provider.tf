# GCP
provider "google" {
  version = "~> 2.18.0"
}

provider "google-beta" {
  version = "~> 2.18.0"
}

# Kubernetes
data "google_client_config" "default" {
}

provider "kubernetes" {
  host                   = module.gke.cluster_endpoint
  cluster_ca_certificate = base64decode(module.gke.cluster_ca)
  token                  = data.google_client_config.default.access_token
  load_config_file       = false
}
