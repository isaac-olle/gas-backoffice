resource "digitalocean_app" "backoffice" {
  spec {
    name   = var.app_name
    region = var.app_region

    service {
      name               = "backoffice"
      instance_count     = var.instance_count
      instance_size_slug = var.instance_size_slug
      http_port          = var.http_port

      image {
        registry_type = "DOCR"
        registry      = data.digitalocean_container_registry.gas.name
        repository    = var.repository
        tag           = var.image_tag

        deploy_on_push {
          enabled = false
        }
      }
    }
  }
}

resource "digitalocean_database_firewall" "gas_postgres" {
  cluster_id = data.digitalocean_database_cluster.postgres.id

  rule {
    type  = "app"
    value = digitalocean_app.backoffice.id
  }
}

resource "digitalocean_project_resources" "backoffice" {
  project   = data.digitalocean_project.main.id
  resources = [digitalocean_app.backoffice.urn]
}
