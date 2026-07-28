data "digitalocean_database_user" "gas_backoffice_service" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice_service"
}

data "digitalocean_database_user" "gas_backoffice_admin" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice_admin"
}

locals {
  postgres_dsn   = "postgresql://${data.digitalocean_database_user.gas_backoffice_service.name}:${data.digitalocean_database_user.gas_backoffice_service.password}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
  migrations_dsn = "postgresql://${data.digitalocean_database_user.gas_backoffice_admin.name}:${data.digitalocean_database_user.gas_backoffice_admin.password}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"

}

resource "digitalocean_database_db" "gas_backoffice" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice"
}

resource "postgresql_schema" "public" {
  name     = "public"
  database = digitalocean_database_db.gas_backoffice.name
  owner    = "doadmin"

  lifecycle {
    ignore_changes = [owner]
  }
}

resource "postgresql_grant" "gas_backoffice_admin_database_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = "gas_backoffice_admin"
  object_type = "database"
  privileges  = ["ALL"]
}

resource "postgresql_grant" "gas_backoffice_service_database_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = "gas_backoffice_service"
  object_type = "database"
  privileges  = ["CONNECT"]
}
