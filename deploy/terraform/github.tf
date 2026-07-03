resource "github_branch" "environment" {
  repository = var.github_repo
  branch     = var.env
}

resource "github_repository_environment" "environment" {
  repository  = var.github_repo
  environment = var.env
}

resource "github_repository_file" "pipeline_trigger" {
  repository          = var.github_repo
  branch              = github_branch.environment.branch
  file                = ".github/workflows/pipeline-${var.env}.yml"
  overwrite_on_create = true

  content = <<-YAML
    name: CI/CD - ${var.env}

    on:
      push:
        branches: [${var.env}]
      pull_request:
        branches: [${var.env}]

    jobs:
      pipeline:
        uses: ./.github/workflows/_pipeline.yml
        with:
          environment: ${var.env}
        secrets: inherit
  YAML
}

resource "github_actions_secret" "do_token" {
  repository      = var.github_repo
  secret_name     = "DO_TOKEN"
  plaintext_value = var.do_token
}

resource "github_actions_secret" "do_app_id" {
  repository      = var.github_repo
  secret_name     = "DO_APP_ID"
  plaintext_value = digitalocean_app.backoffice.id
}

resource "github_actions_secret" "do_registry" {
  repository      = var.github_repo
  secret_name     = "DO_REGISTRY"
  plaintext_value = "registry.digitalocean.com/${var.registry_name}"
}
