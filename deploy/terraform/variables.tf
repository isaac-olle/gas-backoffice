variable "env" { type = string }

variable "postgres_cluster_name" { type = string }
variable "registry_name" { type = string }
variable "project_name" { type = string }
variable "app_name" { type = string }
variable "app_region" { type = string }
variable "instance_count" { type = number }
variable "instance_size_slug" { type = string }
variable "http_port" { type = number }
variable "repository" { type = string }
variable "image_tag" { type = string }
variable "github_repo" { type = string }
variable "source_branch" {
  type    = string
  default = "main"
}

variable "allowed_origins" { type = string }
variable "cookie_origins" { type = string }
variable "seed_token" {
  type      = string
  sensitive = true
}
variable "backoffice_protected_resource_id" { type = string }
variable "backoffice_openid_scope_id" { type = string }
variable "backoffice_email_scope_id" { type = string }
variable "backoffice_profile_scope_id" { type = string }
variable "backoffice_client_id" { type = string }
variable "backoffice_client_secret" {
  type      = string
  sensitive = true
}
variable "backoffice_client_redirect_uri" { type = string }
variable "gas_protocol_url" { type = string }
variable "encryption_key" {
  type      = string
  sensitive = true
}
