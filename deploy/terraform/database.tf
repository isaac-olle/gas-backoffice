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


resource "postgresql_grant" "gas_backoffice_admin_tables_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_admin.name
  schema      = "public"
  object_type = "table"
  privileges  = ["ALL"]
}

resource "postgresql_grant" "gas_backoffice_admin_sequences_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_admin.name
  schema      = "public"
  object_type = "sequence"
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

resource "postgresql_default_privileges" "gas_backoffice_admin_future_tables_privileges" {
  database    = digitalocean_database_db.gas_backoffice.name
  role        = digitalocean_database_user.gas_backoffice_admin.name
  schema      = "public"
  object_type = "table"
  owner       = "doadmin"
  privileges  = ["ALL"]
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

