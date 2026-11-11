#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

ENV="${ENV:-local}"
ENV_FILE="../.env.${ENV}"
if [[ ! -f "$ENV_FILE" ]]; then
  echo "Missing environment file: .env.${ENV}" >&2
  exit 1
fi

set -a
source "$ENV_FILE"
source google_idp.env
set +a

json_escape() {
	local value="$1"
	value="${value//\\/\\\\}"
	value="${value//\"/\\\"}"
	value="${value//$'\b'/\\b}"
	value="${value//$'\f'/\\f}"
	value="${value//$'\n'/\\n}"
	value="${value//$'\r'/\\r}"
	value="${value//$'\t'/\\t}"
	value="${value//&/\\u0026}"
	printf '%s' "$value"
}

payload="$(<seed.json)"
replace_seed_value() {
	local name="$1"
	local marker="\${$name}"
	local value="${!name}"
	value="$(json_escape "$value")"
	payload="${payload//"$marker"/$value}"
}

for name in \
	SEED_ADMIN_NAME SEED_ADMIN_EMAIL SEED_SUBSCRIPTION_NAME \
	GAS_BACKOFFICE_PROTECTED_RESOURCE_ID GAS_BACKOFFICE_OPENID_SCOPE_ID \
	GAS_BACKOFFICE_EMAIL_SCOPE_ID GAS_BACKOFFICE_PROFILE_SCOPE_ID \
	GAS_BACKOFFICE_CLIENT_ID GAS_BACKOFFICE_CLIENT_REDIRECT_URI GAS_BACKOFFICE_CLIENT_SECRET \
	GOOGLE_IDP_ID GOOGLE_ISSUER GOOGLE_CLIENT_ID GOOGLE_CLIENT_SECRET \
	GOOGLE_AUTHORIZATION_ENDPOINT GOOGLE_TOKEN_ENDPOINT GOOGLE_USERINFO_ENDPOINT \
	GOOGLE_JWKS_URI BACKOFFICE_REDIRECT_URI; do
	replace_seed_value "$name"
done

payload_file="$(mktemp)"
response_file="$(mktemp)"
trap 'rm -f "$payload_file" "$response_file"' EXIT
printf '%s\n' "$payload" > "$payload_file"

status="$(curl --silent --show-error \
  --output "$response_file" \
  --write-out '%{http_code}' \
  --request POST \
  --header "Authorization: Bearer $SEED_TOKEN" \
  --header 'Content-Type: application/json' \
  --data-binary "@$payload_file" \
  "${API_BASE_URL%/}/v1/seed")"

if [[ "$status" != "201" ]]; then
  echo "Seed failed with HTTP $status:" >&2
  cat "$response_file" >&2
  exit 1
fi

cat "$response_file"
