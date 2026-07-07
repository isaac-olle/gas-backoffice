resource "postgresql_script" "init_backoffice" {
  database = digitalocean_database_db.gas_backoffice.name
  scripts  = [file("${path.module}/../sql/01-backoffice.sql")]

  lifecycle {
    ignore_changes = all
  }

  depends_on = [
    postgresql_grant.gas_backoffice_admin_database_privileges,
    postgresql_grant.gas_backoffice_admin_schema_privileges,
  ]
}

resource "postgresql_script" "init_protocol" {
  database = digitalocean_database_db.gas_backoffice.name
  scripts  = [file("${path.module}/../sql/02-protocol.sql")]

  lifecycle {
    ignore_changes = all
  }

  depends_on = [postgresql_script.init_backoffice]
}

resource "postgresql_script" "init_security" {
  database = digitalocean_database_db.gas_backoffice.name
  scripts  = [file("${path.module}/../sql/03-security.sql")]

  lifecycle {
    ignore_changes = all
  }

  depends_on = [postgresql_script.init_protocol]
}

resource "postgresql_script" "init_audit_events" {
  database = digitalocean_database_db.gas_backoffice.name
  scripts  = [file("${path.module}/../sql/XX-audit_events.sql")]

  lifecycle {
    ignore_changes = all
  }

  depends_on = [postgresql_script.init_security]
}

resource "postgresql_script" "init_outbox" {
  database = digitalocean_database_db.gas_backoffice.name
  scripts  = [file("${path.module}/../sql/XX-outbox_event_bus.sql")]

  lifecycle {
    ignore_changes = all
  }

  depends_on = [postgresql_script.init_audit_events]
}
