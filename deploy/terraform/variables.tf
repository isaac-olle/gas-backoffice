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
  type = string
}

# GitHub
variable "github_repo" {
  type = string
}

variable "source_branch" {
  type    = string
  default = "main"
}

# Runtime configuration from config/config.yaml
variable "allowed_origins" {
  type    = string
  default = ""
}

variable "cookie_origins" {
  type    = string
  default = ""
}

variable "seed_token" {
  type      = string
  sensitive = true
  default   = ""
}

variable "backoffice_protected_resource_id" {
  type    = string
  default = ""
}

variable "backoffice_openid_scope_id" {
  type    = string
  default = ""
}

variable "backoffice_email_scope_id" {
  type    = string
  default = ""
}

variable "backoffice_profile_scope_id" {
  type    = string
  default = ""
}

variable "backoffice_client_id" {
  type    = string
  default = ""
}

variable "backoffice_client_secret" {
  type      = string
  sensitive = true
  default   = ""
}

variable "backoffice_client_redirect_uri" {
  type    = string
  default = ""
}

variable "gas_protocol_url" {
  type    = string
  default = ""
}

variable "encryption_key" {
  type      = string
  sensitive = true
  default   = ""
}
