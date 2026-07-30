resource "github_actions_variable" "do_app_id" {
  repository    = var.github_repo
  variable_name = "${var.env}_DO_APP_ID"
  value         = digitalocean_app.backoffice.id
}
