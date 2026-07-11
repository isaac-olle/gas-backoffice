resource "digitalocean_app" "backoffice" {
  spec {
    name   = var.app_name
    region = var.app_region

    service {
      name               = "backoffice"
      instance_count     = var.instance_count
      instance_size_slug = var.instance_size_slug
      http_port          = var.http_port

      health_check {
        http_path             = "/health"
        initial_delay_seconds = 10
        period_seconds        = 30
        timeout_seconds       = 5
        success_threshold     = 1
        failure_threshold     = 3
        port                  = var.http_port
      }

      env {
        key   = "ENV"
        value = var.env
      }

      env {
        key   = "API_HOST"
        value = "0.0.0.0"
      }

      env {
        key   = "API_PORT"
        value = tostring(var.http_port)
      }

      dynamic "env" {
        for_each = var.backoffice_pr_id != "" ? [1] : []
        content {
          key   = "BACKOFFICE_PR_ID"
          value = var.backoffice_pr_id
          type  = "SECRET"
        }
      }

      env {
        key   = "POSTGRES_DSN"
        value = "postgresql://${digitalocean_database_user.gas_backoffice_service.name}:${digitalocean_database_user.gas_backoffice_service.password}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
        type  = "SECRET"
      }

      env {
        key   = "MIGRATIONS_DSN"
        value = "postgresql://${digitalocean_database_user.gas_backoffice_admin.name}:${digitalocean_database_user.gas_backoffice_admin.password}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
        type  = "SECRET"
      }

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

  lifecycle {
    ignore_changes = [spec[0].service[0].image[0].tag]
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
