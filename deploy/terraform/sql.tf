# These PostgreSQL-managed users are kept here for reference. The active users
# are managed by DigitalOcean below.
# resource "random_password" "gas_backoffice_admin" {
#   length           = 32
#   special          = true
#   override_special = "!#$%&*()-_=+"
# }
#
# resource "random_password" "gas_backoffice_service" {
#   length           = 32
#   special          = true
#   override_special = "!#$%&*()-_=+"
# }
#
# resource "postgresql_role" "gas_backoffice_admin" {
#   name           = "gas_backoffice_admin"
#   login          = true
#   password       = random_password.gas_backoffice_admin.result
#   skip_drop_role = true
# }
#
# resource "postgresql_role" "gas_backoffice_service" {
#   name     = "gas_backoffice_service"
#   login    = true
#   password = random_password.gas_backoffice_service.result
# }

# Manage cluster users through DigitalOcean, which also provisions their passwords.
resource "digitalocean_database_user" "gas_backoffice_admin" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice_admin"
}

resource "digitalocean_database_user" "gas_backoffice_service" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice_service"
}

resource "digitalocean_database_db" "gas_backoffice" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice"
}

locals {
  postgres_dsn   = "postgresql://${digitalocean_database_user.gas_backoffice_service.name}:${urlencode(digitalocean_database_user.gas_backoffice_service.password)}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
  migrations_dsn = "postgresql://${digitalocean_database_user.gas_backoffice_admin.name}:${urlencode(digitalocean_database_user.gas_backoffice_admin.password)}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
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
  privileges  = ["USAGE"]
}
