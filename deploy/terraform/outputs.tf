output "app_id" {
  value = digitalocean_app.backoffice.id
}

output "app_url" {
  value = digitalocean_app.backoffice.live_url
}

output "db_user_admin" {
  value = digitalocean_database_user.gas_admin.name
}

output "db_user_backoffice" {
  value = digitalocean_database_user.gas_backoffice.name
}
