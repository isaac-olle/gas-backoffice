data "digitalocean_database_cluster" "postgres" {
  name = var.postgres_cluster_name
}

data "digitalocean_container_registry" "gas" {
  name = var.registry_name
}

data "digitalocean_project" "main" {
  name = var.project_name
}
