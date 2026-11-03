#!/usr/bin/env bash
set -euo pipefail

env_name="${1:-}"
if [[ -z "$env_name" || ! "$env_name" =~ ^[A-Za-z0-9_-]+$ ]]; then
	echo 'Usage: make env ENV=dev' >&2
	exit 1
fi

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
env_file="$project_root/.env.$env_name"
if [[ -e "$env_file" ]]; then
	echo "Environment file already exists: .env.$env_name" >&2
	exit 1
fi

cp "$project_root/.env.example" "$env_file"
echo "Created .env.$env_name from .env.example; edit it with your environment values."
