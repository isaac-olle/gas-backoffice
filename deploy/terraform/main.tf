# ---- Data Sources ----

data "digitalocean_database_cluster" "postgres" {
  name = var.postgres_cluster_name
}

data "digitalocean_container_registry" "gas" {
  name = var.registry_name
}

data "digitalocean_project" "main" {
  name = var.project_name
}

# ---- Database ----

resource "digitalocean_database_db" "gas_backoffice" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice"
}

resource "digitalocean_database_user" "gas_backoffice_admin" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice_admin"
}

resource "digitalocean_database_user" "gas_backoffice_service" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice_service"
}

resource "postgresql_grant" "gas_backoffice_admin_database_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_admin.name
  object_type = "database"
  privileges  = ["ALL"]
}

resource "postgresql_grant" "gas_backoffice_admin_schema_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_admin.name
  schema      = "public"
  object_type = "schema"
  privileges  = ["ALL"]
}

resource "postgresql_grant" "gas_backoffice_service_database_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_service.name
  object_type = "database"
  privileges  = ["CONNECT"]
}

resource "postgresql_grant" "gas_backoffice_service_schema_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_service.name
  schema      = "public"
  object_type = "schema"
  privileges  = ["USAGE"]
}

resource "postgresql_grant" "gas_backoffice_service_tables_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_service.name
  schema      = "public"
  object_type = "table"
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE"]
}

resource "postgresql_grant" "gas_backoffice_service_sequences_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_service.name
  schema      = "public"
  object_type = "sequence"
  privileges  = ["USAGE", "SELECT"]
}

resource "postgresql_default_privileges" "gas_backoffice_service_future_tables_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_service.name
  schema      = "public"
  object_type = "table"
  owner       = digitalocean_database_user.gas_backoffice_admin.name
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE"]
}

resource "postgresql_default_privileges" "gas_backoffice_service_future_sequences_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_service.name
  schema      = "public"
  object_type = "sequence"
  owner       = digitalocean_database_user.gas_backoffice_admin.name
  privileges  = ["USAGE", "SELECT"]
}

# ---- App Platform ----

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

# ---- GitHub Branch & Environment ----

resource "github_branch" "environment" {
  repository = var.github_repo
  branch     = var.env
}

resource "github_repository_environment" "environment" {
  repository  = var.github_repo
  environment = var.env
}

resource "github_repository_file" "pipeline_trigger" {
  repository          = var.github_repo
  branch              = github_branch.environment.branch
  file                = ".github/workflows/pipeline-${var.env}.yml"
  overwrite_on_create = true

  content = <<-YAML
    name: CI/CD - ${var.env}

    on:
      push:
        branches: [${var.env}]
      pull_request:
        branches: [${var.env}]

    jobs:
      pipeline:
        uses: ./.github/workflows/_pipeline.yml
        with:
          environment: ${var.env}
        secrets: inherit
  YAML
}

# ---- GitHub Secrets ----

resource "github_actions_secret" "do_token" {
  repository      = var.github_repo
  secret_name     = "DO_TOKEN"
  plaintext_value = var.do_token
}

resource "github_actions_secret" "do_app_id" {
  repository      = var.github_repo
  secret_name     = "DO_APP_ID"
  plaintext_value = digitalocean_app.backoffice.id
}

resource "github_actions_secret" "do_registry" {
  repository      = var.github_repo
  secret_name     = "DO_REGISTRY"
  plaintext_value = "registry.digitalocean.com/${data.digitalocean_container_registry.gas.name}"
}
