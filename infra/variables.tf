variable "region" {
  type    = string
  default = "lon1"
}

variable "vpc_range" {
  type = string
  default = "10.0.0.0/24"
}

variable "doks_version" {
  type = string
  default = "1.29.1-do.0"
}

variable "doks_node_pool_size" {
  type = string
  default = "s-2vcpu-2gb"
}

variable "doks_node_count" {
  type = number
  default = 3
}

variable "name" {
  type = string
}