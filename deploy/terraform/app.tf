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
        key   = "ALLOWED_ORIGINS"
        value = var.allowed_origins
      }

      env {
        key   = "COOKIE_ORIGINS"
        value = var.cookie_origins
      }

      env {
        key   = "SEED_TOKEN"
        value = var.seed_token
        type  = "SECRET"
      }

      env {
        key   = "GAS_BACKOFFICE_PROTECTED_RESOURCE_ID"
        value = var.backoffice_protected_resource_id
      }

      env {
        key   = "GAS_BACKOFFICE_OPENID_SCOPE_ID"
        value = var.backoffice_openid_scope_id
      }

      env {
        key   = "GAS_BACKOFFICE_EMAIL_SCOPE_ID"
        value = var.backoffice_email_scope_id
      }

      env {
        key   = "GAS_BACKOFFICE_PROFILE_SCOPE_ID"
        value = var.backoffice_profile_scope_id
      }

      env {
        key   = "GAS_BACKOFFICE_CLIENT_ID"
        value = var.backoffice_client_id
      }

      env {
        key   = "GAS_BACKOFFICE_CLIENT_SECRET"
        value = var.backoffice_client_secret
        type  = "SECRET"
      }

      env {
        key   = "GAS_BACKOFFICE_CLIENT_REDIRECT_URI"
        value = var.backoffice_client_redirect_uri
      }

      env {
        key   = "GAS_PROTOCOL_URL"
        value = var.gas_protocol_url
      }

      env {
        key   = "ENCRYPTION_KEY"
        value = var.encryption_key
        type  = "SECRET"
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
      name               = "pre-deploy"
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
        key   = "POSTGRES_DSN"
        value = local.migrations_dsn
        type  = "SECRET"
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
  }

  lifecycle {
    ignore_changes = [
      spec[0].service[0].image[0].tag,
      spec[0].job[0].image[0].tag,
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
