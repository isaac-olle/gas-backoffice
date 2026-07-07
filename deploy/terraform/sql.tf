locals {
  postgres_admin_dsn = "postgresql://${data.digitalocean_database_cluster.postgres.user}:${data.digitalocean_database_cluster.postgres.password}@${data.digitalocean_database_cluster.postgres.host}:${data.digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.gas_backoffice.name}?sslmode=require"
}

resource "null_resource" "init_sql" {
  provisioner "local-exec" {
    command = <<-EOT
      psql "${local.postgres_admin_dsn}" -f "${path.module}/../sql/01-backoffice.sql"
      psql "${local.postgres_admin_dsn}" -f "${path.module}/../sql/02-protocol.sql"
      psql "${local.postgres_admin_dsn}" -f "${path.module}/../sql/03-security.sql"
      psql "${local.postgres_admin_dsn}" -f "${path.module}/../sql/XX-audit_events.sql"
      psql "${local.postgres_admin_dsn}" -f "${path.module}/../sql/XX-outbox_event_bus.sql"
    EOT
  }

  lifecycle {
    ignore_changes = all
  }

  depends_on = [
    postgresql_grant.gas_backoffice_admin_database_privileges,
    postgresql_grant.gas_backoffice_admin_schema_privileges,
  ]
}
