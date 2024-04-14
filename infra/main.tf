terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
  backend "s3" {
    endpoints                    = {
      s3 = "https://fra1.digitaloceanspaces.com"
    }
    key                         = "terraform.tfstate"
    bucket                      = "dscott-tfstate"
    region                      = "us-west-1"
    skip_requesting_account_id  = true
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_metadata_api_check     = true
    skip_s3_checksum            = true
    use_path_style              = true
  }
}

provider "digitalocean" {}

resource "digitalocean_vpc" "vpc" {
  name     = var.name
  region   = var.region
  ip_range = var.vpc_range
}

resource "digitalocean_kubernetes_cluster" "cluster" {
  name   = var.name
  region = var.region
  version = var.doks_version
  vpc_uuid = digitalocean_vpc.vpc.id

  node_pool {
    name       = "${var.name}-pool"
    size       = var.doks_node_pool_size
    node_count = var.doks_node_count
  }
}