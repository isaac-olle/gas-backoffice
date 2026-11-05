resource "digitalocean_database_user" "gas_backoffice_service" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas-backoffice-user"

  # DigitalOcean's API returns an empty settings object for PostgreSQL users.
  # Keep it configured to avoid the provider's phantom settings update.
  settings {}

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_database_db" "gas_backoffice" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice"
}

data "digitalocean_database_user" "gas_admin" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_admin"
}

locals {
  postgres_dsn   = "postgresql://${digitalocean_database_user.gas_backoffice_service.name}:${urlencode(digitalocean_database_user.gas_backoffice_service.password)}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
  # gas_admin is created by the parent infrastructure project and is used for
  # schema migrations; the application itself uses the restricted service user.
  migrations_dsn = "postgresql://${data.digitalocean_database_user.gas_admin.name}:${urlencode(data.digitalocean_database_user.gas_admin.password)}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
}
