resource "null_resource" "init_sql" {
  provisioner "local-exec" {
    command = <<-EOT
      psql "${var.postgres_admin_dsn}" -f "${path.module}/../sql/01-backoffice.sql"
      psql "${var.postgres_admin_dsn}" -f "${path.module}/../sql/02-protocol.sql"
      psql "${var.postgres_admin_dsn}" -f "${path.module}/../sql/03-security.sql"
      psql "${var.postgres_admin_dsn}" -f "${path.module}/../sql/XX-audit_events.sql"
      psql "${var.postgres_admin_dsn}" -f "${path.module}/../sql/XX-outbox_event_bus.sql"
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
