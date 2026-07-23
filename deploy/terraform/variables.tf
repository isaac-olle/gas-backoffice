# Environment
variable "env" {
  type = string
}

# Data source lookups
variable "postgres_cluster_name" {
  type = string
}

variable "registry_name" {
  type = string
}

variable "project_name" {
  type = string
}

# App Platform
variable "app_name" {
  type = string
}

variable "app_region" {
  type = string
}

variable "instance_count" {
  type = number
}

variable "instance_size_slug" {
  type = string
}

variable "http_port" {
  type = number
}

variable "repository" {
  type = string
}

variable "image_tag" {
  type    = string
  default = "dev"
}

# App environment
variable "gas_seed_admin_name" {
  type    = string
  default = "Gas-Seed-Admin"
}

variable "gas_backoffice_pr_name" {
  type    = string
  default = "Gas-Backoffice"
}

# GitHub
variable "github_repo" {
  type = string
}

variable "source_branch" {
  type    = string
  default = "main"
}
