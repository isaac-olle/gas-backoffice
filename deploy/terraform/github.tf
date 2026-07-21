resource "github_branch" "environment" {
  repository    = var.github_repo
  branch        = var.env
  source_branch = var.source_branch

  lifecycle {
    prevent_destroy = true
  }
}

resource "github_repository_file" "pipeline_trigger" {
  repository          = var.github_repo
  branch              = github_branch.environment.branch
  file                = ".github/workflows/pipeline-${var.env}.yml"
  overwrite_on_create = true

  lifecycle {
    prevent_destroy = true
  }

  content = <<-YAML
    name: CI/CD - ${var.env}

    on:
      push:
        branches: [${var.env}]
      pull_request:
        branches: [${var.env}]
      workflow_dispatch:

    jobs:
      pipeline:
        uses: ./.github/workflows/_pipeline.yml
        with:
          environment: ${var.env}
        secrets: inherit
  YAML
}

resource "github_actions_secret" "do_app_id" {
  repository      = var.github_repo
  secret_name     = "${upper(var.env)}_DO_APP_ID"
  plaintext_value = digitalocean_app.backoffice.id
}
