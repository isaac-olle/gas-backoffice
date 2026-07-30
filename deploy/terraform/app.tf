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

      env {
        key   = "GAS_SEED_ADMIN_NAME"
        value = "Gas Seed Admin"
      }

      env {
        key   = "GAS_BACKOFFICE_PR_NAME"
        value = "Gas Backoffice"
      }

      env {
        key   = "POSTGRES_DSN"
        value = local.postgres_dsn
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

    job {
      name               = "migrate"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb"
      kind               = "PRE_DEPLOY"
      run_command        = "/migrate"

      env {
        key   = "MIGRATIONS_DSN"
        value = local.migrations_dsn
        type  = "SECRET"
      }

      env {
        key   = "MIGRATIONS_PATH"
        value = "/golang-migrate"
      }

      env {
        key   = "RUN_MIGRATIONS"
        value = "true"
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

    job {
      name               = "db-seed"
      instance_count     = 1
      instance_size_slug = "apps-s-1vcpu-0.5gb"
      kind               = "PRE_DEPLOY"
      run_command        = "/db-seed"

      env {
        key   = "POSTGRES_DSN"
        value = local.migrations_dsn
        type  = "SECRET"
      }

      env {
        key   = "GAS_SEED_ADMIN_NAME"
        value = var.gas_seed_admin_name
      }

      env {
        key   = "GAS_BACKOFFICE_PR_NAME"
        value = var.gas_backoffice_pr_name
      }

      env {
        key   = "RUN_SEED"
        value = "true"
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
    ignore_changes = [
      spec[0].service[0].image[0].tag,
      spec[0].job[0].image[0].tag,
      spec[0].job[1].image[0].tag,
    ]
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

output "app_id" {
  value = digitalocean_app.backoffice.id
}
