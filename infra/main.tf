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
  name     = "capi-vpc-demo"
  region   = "lon1"
  ip_range = "10.0.0.0/24"
}