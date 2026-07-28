resource "random_password" "gas_backoffice_admin" {
  length            = 32
  special           = true
  override_special  = "!#$%&*()-_=+"
}

resource "random_password" "gas_backoffice_service" {
  length            = 32
  special           = true
  override_special  = "!#$%&*()-_=+"
}

resource "postgresql_role" "gas_backoffice_admin" {
  name     = "gas_backoffice_admin"
  login    = true
  password = random_password.gas_backoffice_admin.result
}

resource "postgresql_role" "gas_backoffice_service" {
  name     = "gas_backoffice_service"
  login    = true
  password = random_password.gas_backoffice_service.result
}

locals {
  postgres_dsn   = "postgresql://gas_backoffice_service:${urlencode(random_password.gas_backoffice_service.result)}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
  migrations_dsn = "postgresql://gas_backoffice_admin:${urlencode(random_password.gas_backoffice_admin.result)}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
}

resource "digitalocean_database_db" "gas_backoffice" {
  cluster_id = data.digitalocean_database_cluster.postgres.id
  name       = "gas_backoffice"
}


resource "postgresql_grant" "gas_backoffice_admin_database_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = postgresql_role.gas_backoffice_admin.name
  object_type = "database"
  privileges  = ["ALL"]
}

resource "postgresql_grant" "gas_backoffice_service_database_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = postgresql_role.gas_backoffice_service.name
  object_type = "database"
  privileges  = ["CONNECT"]
}

resource "postgresql_grant" "gas_backoffice_service_tables_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = postgresql_role.gas_backoffice_service.name
  schema      = "public"
  object_type = "table"
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE"]
}

resource "postgresql_grant" "gas_backoffice_service_sequences_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = postgresql_role.gas_backoffice_service.name
  schema      = "public"
  object_type = "sequence"
  privileges  = ["USAGE"]
}
