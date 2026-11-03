#!/usr/bin/env bash
set -euo pipefail

env_name="${1:-}"
action="${2:-plan}"
if [[ -z "$env_name" || ! "$env_name" =~ ^[A-Za-z0-9_-]+$ ]]; then
	echo 'Usage: make terraform-plan ENV=dev or make terraform-apply ENV=dev' >&2
	exit 1
fi
if [[ "$action" != "plan" && "$action" != "apply" ]]; then
	echo 'Terraform action must be plan or apply' >&2
	exit 1
fi

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
env_file="$project_root/.env.$env_name"
if [[ ! -f "$env_file" ]]; then
	echo "Environment file not found: $env_file (create it with 'make env ENV=$env_name')" >&2
	exit 1
fi

# Load dotenv values into this process without printing secrets.
set -a
# shellcheck disable=SC1090
source "$env_file"
set +a

export TF_VAR_env="${ENV:-$env_name}"
export TF_VAR_postgres_cluster_name="${POSTGRES_CLUSTER_NAME:-}"
export TF_VAR_registry_name="${REGISTRY_NAME:-}"
export TF_VAR_project_name="${PROJECT_NAME:-}"
export TF_VAR_app_name="${APP_NAME:-}"
export TF_VAR_app_region="${APP_REGION:-}"
export TF_VAR_instance_count="${INSTANCE_COUNT:-1}"
export TF_VAR_instance_size_slug="${INSTANCE_SIZE_SLUG:-apps-s-1vcpu-0.5gb}"
export TF_VAR_http_port="${API_PORT:-8090}"
export TF_VAR_repository="${REPOSITORY:-}"
export TF_VAR_image_tag="${IMAGE_TAG:-latest}"
export TF_VAR_github_repo="${GITHUB_REPO:-}"
export TF_VAR_source_branch="${SOURCE_BRANCH:-main}"
export TF_VAR_allowed_origins="${ALLOWED_ORIGINS:-}"
export TF_VAR_cookie_origins="${COOKIE_ORIGINS:-}"
export TF_VAR_seed_token="${SEED_TOKEN:-}"
export TF_VAR_backoffice_protected_resource_id="${GAS_BACKOFFICE_PROTECTED_RESOURCE_ID:-}"
export TF_VAR_backoffice_openid_scope_id="${GAS_BACKOFFICE_OPENID_SCOPE_ID:-}"
export TF_VAR_backoffice_email_scope_id="${GAS_BACKOFFICE_EMAIL_SCOPE_ID:-}"
export TF_VAR_backoffice_profile_scope_id="${GAS_BACKOFFICE_PROFILE_SCOPE_ID:-}"
export TF_VAR_backoffice_client_id="${GAS_BACKOFFICE_CLIENT_ID:-}"
export TF_VAR_backoffice_client_secret="${GAS_BACKOFFICE_CLIENT_SECRET:-}"
export TF_VAR_backoffice_client_redirect_uri="${GAS_BACKOFFICE_CLIENT_REDIRECT_URI:-}"
export TF_VAR_gas_protocol_url="${GAS_PROTOCOL_URL:-}"
export TF_VAR_encryption_key="${ENCRYPTION_KEY:-}"

terraform_dir="$project_root/deploy/terraform"
terraform -chdir="$terraform_dir" init -input=false
terraform -chdir="$terraform_dir" "$action" -input=false
