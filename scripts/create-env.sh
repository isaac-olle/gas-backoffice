#!/usr/bin/env bash
set -euo pipefail

env_name="${1:-}"
if [[ -z "$env_name" || ! "$env_name" =~ ^[A-Za-z0-9_-]+$ ]]; then
	echo 'Usage: make env ENV=dev' >&2
	exit 1
fi

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
env_file="$project_root/.env.$env_name"
template_file="$project_root/.env.example"
if [[ -e "$env_file" ]]; then
	echo "Environment file already exists: .env.$env_name" >&2
	exit 1
fi
if [[ ! -f "$template_file" ]]; then
	echo "Environment template not found: $template_file" >&2
	exit 1
fi
needs_uuidgen=false
needs_openssl=false
command -v uuidgen >/dev/null 2>&1 || needs_uuidgen=true
command -v openssl >/dev/null 2>&1 || needs_openssl=true

if [[ "$needs_uuidgen" == true || "$needs_openssl" == true ]]; then
	install_packages() {
		if [[ "$EUID" -eq 0 ]]; then
			"$@"
		elif command -v sudo >/dev/null 2>&1; then
			sudo "$@"
		else
			echo 'Installing required packages needs root privileges or sudo.' >&2
			exit 1
		fi
	}

	if command -v apt-get >/dev/null 2>&1; then
		packages=()
		[[ "$needs_uuidgen" == false ]] || packages+=(uuid-runtime)
		[[ "$needs_openssl" == false ]] || packages+=(openssl)
		echo "Installing required packages: ${packages[*]}"
		install_packages apt-get update
		install_packages apt-get install -y "${packages[@]}"
	elif command -v apk >/dev/null 2>&1; then
		packages=()
		[[ "$needs_uuidgen" == false ]] || packages+=(uuidgen)
		[[ "$needs_openssl" == false ]] || packages+=(openssl)
		echo "Installing required packages: ${packages[*]}"
		install_packages apk add --no-cache "${packages[@]}"
	elif command -v dnf >/dev/null 2>&1; then
		packages=()
		[[ "$needs_uuidgen" == false ]] || packages+=(util-linux)
		[[ "$needs_openssl" == false ]] || packages+=(openssl)
		echo "Installing required packages: ${packages[*]}"
		install_packages dnf install -y "${packages[@]}"
	elif command -v yum >/dev/null 2>&1; then
		packages=()
		[[ "$needs_uuidgen" == false ]] || packages+=(util-linux)
		[[ "$needs_openssl" == false ]] || packages+=(openssl)
		echo "Installing required packages: ${packages[*]}"
		install_packages yum install -y "${packages[@]}"
	elif command -v pacman >/dev/null 2>&1; then
		packages=()
		[[ "$needs_uuidgen" == false ]] || packages+=(util-linux)
		[[ "$needs_openssl" == false ]] || packages+=(openssl)
		echo "Installing required packages: ${packages[*]}"
		install_packages pacman -Sy --noconfirm "${packages[@]}"
	elif command -v zypper >/dev/null 2>&1; then
		packages=()
		[[ "$needs_uuidgen" == false ]] || packages+=(util-linux)
		[[ "$needs_openssl" == false ]] || packages+=(openssl)
		echo "Installing required packages: ${packages[*]}"
		install_packages zypper --non-interactive install "${packages[@]}"
	else
		echo 'No supported Linux package manager found to install uuidgen and openssl.' >&2
		exit 1
	fi
fi

if ! command -v uuidgen >/dev/null 2>&1 || ! command -v openssl >/dev/null 2>&1; then
	still_missing=()
	command -v uuidgen >/dev/null 2>&1 || still_missing+=(uuidgen)
	command -v openssl >/dev/null 2>&1 || still_missing+=(openssl)
	echo "Package installation completed, but these commands are still unavailable: ${still_missing[*]}" >&2
	exit 1
fi

protected_resource_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"
openid_scope_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"
email_scope_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"
profile_scope_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"
client_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"
encryption_key="$(openssl rand -base64 32)"
seed_token="$(openssl rand -hex 32)"
client_secret="$(openssl rand -hex 32)"

awk \
	-v env_name="$env_name" \
	-v protected_resource_id="$protected_resource_id" \
	-v openid_scope_id="$openid_scope_id" \
	-v email_scope_id="$email_scope_id" \
	-v profile_scope_id="$profile_scope_id" \
	-v client_id="$client_id" \
	-v client_secret="$client_secret" \
	-v encryption_key="$encryption_key" \
	-v seed_token="$seed_token" '
BEGIN {
	values["ENV"] = env_name
	values["GAS_BACKOFFICE_PROTECTED_RESOURCE_ID"] = protected_resource_id
	values["GAS_BACKOFFICE_OPENID_SCOPE_ID"] = openid_scope_id
	values["GAS_BACKOFFICE_EMAIL_SCOPE_ID"] = email_scope_id
	values["GAS_BACKOFFICE_PROFILE_SCOPE_ID"] = profile_scope_id
	values["GAS_BACKOFFICE_CLIENT_ID"] = client_id
	values["GAS_BACKOFFICE_CLIENT_SECRET"] = client_secret
	values["ENCRYPTION_KEY"] = encryption_key
	values["SEED_TOKEN"] = seed_token
}
{
	separator = index($0, "=")
	if (separator > 0) {
		key = substr($0, 1, separator - 1)
		if (key in values) {
			printf "%s=\"%s\"\n", key, values[key]
			next
		}
	}
	print
}' "$template_file" > "$env_file"

echo "Created .env.$env_name with generated IDs, an AES-256 key, and a seed token. Review and set the remaining environment values."
