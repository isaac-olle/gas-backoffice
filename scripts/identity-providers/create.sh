#!/usr/bin/env bash
set +x
set -euo pipefail

usage() {
  cat <<'EOF'
Create an OpenID provider through the GAS admin API (server-side AES encryption).

Microsoft:
  bash scripts/identity-providers/create.sh \
    --api-base-url https://gas.example.com --tenant-id TENANT_GUID --client-id CLIENT_ID

Other providers:
  bash scripts/identity-providers/create.sh \
    --api-base-url https://gas.example.com \
    --discovery-url https://issuer.example/.well-known/openid-configuration \
    --reference-name example --name Example --client-id CLIENT_ID

Options:
  --env-file PATH       Load Bash environment settings before command-line overrides
  --redirect-uri URL    Default: API_BASE_URL/v1/oidc/callback/REFERENCE_NAME
  --reference-name NAME Default: microsoft
  --name NAME           Default: Microsoft
  --scopes "..."        Default: "openid email profile"
  --help

Requires Bash, curl and jq. Prompts silently for the GAS admin access token
and provider secret. Alternatively use GAS_ADMIN_ACCESS_TOKEN and OIDC_CLIENT_SECRET.
API_BASE_URL, OIDC_CLIENT_ID and OIDC_TENANT_ID can supply the corresponding options.
EOF
}

die() { printf '%s\n' "$*" >&2; exit 1; }

env_file=''
arguments=()
while (($#)); do
  if [[ "$1" == '--env-file' ]]; then
    (($# >= 2)) || die 'Missing value for --env-file'
    env_file="$2"
    shift 2
  else
    arguments+=("$1")
    shift
  fi
done
if [[ -n "$env_file" ]]; then
  [[ -f "$env_file" ]] || die "Environment file not found: $env_file"
  # Like seed.sh, this sources a trusted, local Bash environment file.
  source "$env_file"
  set +x
fi
set -- "${arguments[@]}"

api_base_url="${API_BASE_URL:-}"
tenant_id="${OIDC_TENANT_ID:-}"
client_id="${MICROSOFT_CLIENT_ID:-${OIDC_CLIENT_ID:-}}"
discovery_url="${OIDC_DISCOVERY_URL:-}"
reference_name="${OIDC_REFERENCE_NAME:-microsoft}"
display_name="${OIDC_NAME:-Microsoft}"
redirect_uri="${BACKOFFICE_REDIRECT_URI:-${OIDC_REDIRECT_URI:-}}"
scopes="${OIDC_SCOPES:-openid email profile}"

while (($#)); do
  case "$1" in
    --help|-h) usage; exit 0 ;;
    --api-base-url|--tenant-id|--client-id|--discovery-url|--reference-name|--name|--redirect-uri|--scopes)
      (($# >= 2)) || die "Missing value for $1"
      case "$1" in
        --api-base-url) api_base_url="$2" ;;
        --tenant-id) tenant_id="$2" ;;
        --client-id) client_id="$2" ;;
        --discovery-url) discovery_url="$2" ;;
        --reference-name) reference_name="$2" ;;
        --name) display_name="$2" ;;
        --redirect-uri) redirect_uri="$2" ;;
        --scopes) scopes="$2" ;;
      esac
      shift 2 ;;
    *) die "Unknown option: $1 (use --help)" ;;
  esac
done

for dependency in curl jq; do
  command -v "$dependency" >/dev/null 2>&1 || die "Required command not found: $dependency"
done
[[ -n "$api_base_url" && -n "$client_id" ]] || die '--api-base-url and --client-id are required.'
[[ "$reference_name" =~ ^[a-zA-Z0-9_-]+$ ]] || die 'Invalid reference name.'
[[ -n "$display_name" ]] || die '--name must not be empty.'
[[ "$api_base_url" =~ ^https?://[^/?#[:space:]@]+(/[^?#[:space:]]*)?$ ]] || die 'API base URL must be absolute HTTP(S), without user information, query or fragment.'
api_base_url="${api_base_url%/}"

metadata=''
if [[ -n "${MICROSOFT_ISSUER:-}${MICROSOFT_AUTHORIZATION_ENDPOINT:-}${MICROSOFT_TOKEN_ENDPOINT:-}${MICROSOFT_JWKS_URI:-}" ]]; then
  for field in MICROSOFT_ISSUER MICROSOFT_AUTHORIZATION_ENDPOINT MICROSOFT_TOKEN_ENDPOINT MICROSOFT_JWKS_URI; do
    [[ -n "${!field:-}" ]] || die "Missing $field in the Microsoft environment file."
  done
  metadata="$(jq -cn --arg issuer "$MICROSOFT_ISSUER" \
    --arg authorization "$MICROSOFT_AUTHORIZATION_ENDPOINT" --arg token "$MICROSOFT_TOKEN_ENDPOINT" \
    --arg userinfo "${MICROSOFT_USERINFO_ENDPOINT:-}" --arg jwks "$MICROSOFT_JWKS_URI" \
    '{issuer: $issuer, authorization_endpoint: $authorization, token_endpoint: $token, userinfo_endpoint: $userinfo, jwks_uri: $jwks}')"
elif [[ -z "$discovery_url" ]]; then
  [[ "$tenant_id" =~ ^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$ ]] || die 'Provide a concrete Microsoft tenant GUID with --tenant-id, or use --discovery-url.'
  [[ "$tenant_id" != '00000000-0000-0000-0000-000000000000' ]] || die 'Tenant ID must not be empty.'
  discovery_url="https://login.microsoftonline.com/$tenant_id/v2.0/.well-known/openid-configuration"
fi
redirect_uri="${redirect_uri:-$api_base_url/v1/oidc/callback/$reference_name}"

if [[ -z "$metadata" ]]; then
  [[ "$discovery_url" =~ ^https?:// ]] || die 'Discovery URL must use HTTP(S).'
  if ! metadata="$(curl --silent --show-error --fail --location --max-redirs 5 \
    --connect-timeout 10 --max-time 30 "$discovery_url")"; then
    die 'Could not retrieve OpenID Discovery metadata. Check the URL/tenant and connectivity.'
  fi
fi

# Validate metadata before asking for credentials. A templated Microsoft issuer
# is incompatible with GAS's exact issuer comparison.
if ! metadata="$(printf '%s' "$metadata" | jq -ce '
  def http_url: type == "string" and test("^https?://[^/?#[:space:]@]+([/?][^#[:space:]]*)?$");
  if type != "object" then error("Expected a discovery object") else . end |
  . as $m |
  if (["issuer", "authorization_endpoint", "token_endpoint", "jwks_uri"] |
      all(.[]; $m[.] | http_url and (test("[{}]") | not))) and
      (($m.userinfo_endpoint // "") == "" or ($m.userinfo_endpoint | http_url))
  then . else error("Invalid or templated discovery metadata; use a tenant-specific URL") end
')"; then
  die 'Invalid OpenID Discovery metadata; Microsoft requires a concrete tenant-specific issuer.'
fi
printf '%s' "$redirect_uri" | jq -Re 'test("^https?://[^/?#[:space:]@]+([/?][^#[:space:]]*)?$")' >/dev/null || die 'Invalid redirect URI.'
[[ "$scopes" =~ [^[:space:]] ]] || die 'At least one scope is required.'

admin_token="${GAS_ADMIN_ACCESS_TOKEN:-}"
client_secret="${MICROSOFT_CLIENT_SECRET:-${OIDC_CLIENT_SECRET:-}}"
if [[ -z "$admin_token" ]]; then
  read -r -s -p 'GAS system_administrator access token (not SEED_TOKEN): ' admin_token || die 'Unable to read admin token.'
  printf '\n' >&2
fi
if [[ -z "$client_secret" ]]; then
  read -r -s -p 'Provider client secret VALUE (not secret ID): ' client_secret || die 'Unable to read client secret.'
  printf '\n' >&2
fi
admin_token="${admin_token#Bearer }"
[[ -n "$admin_token" && -n "$client_secret" ]] || die 'Admin token and client secret must not be empty.'
[[ "$admin_token" != *$'\n'* && "$admin_token" != *$'\r'* ]] || die 'Admin token must not contain line breaks.'

# JSON encoding preserves quotes, backslashes and other characters in secrets.
# The secret is supplied via the environment rather than a process argument.
payload="$(OIDC_SECRET_VALUE="$client_secret" jq -cn \
  --argjson metadata "$metadata" --arg reference "$reference_name" --arg name "$display_name" \
  --arg client "$client_id" --arg scopes "$scopes" --arg redirect "$redirect_uri" '
  {
    reference_name: $reference, name: $name, issuer: $metadata.issuer,
    client_id: $client, client_secret: env.OIDC_SECRET_VALUE,
    authorization_endpoint: $metadata.authorization_endpoint,
    token_endpoint: $metadata.token_endpoint,
    userinfo_endpoint: ($metadata.userinfo_endpoint // ""), jwks_uri: $metadata.jwks_uri,
    scopes: [$scopes | scan("\\S+")], backoffice_redirect_uri: $redirect
  }')"

config_escape() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  printf '%s' "$value"
}

# Keep the Authorization header and body out of curl's command-line arguments
# and off disk. Do not follow a redirect with this authenticated request.
if ! status="$(
  {
    printf 'url = "%s"\n' "$(config_escape "$api_base_url/v1/backoffice/identity-providers")"
    printf 'request = "POST"\n'
    printf 'header = "Authorization: Bearer %s"\n' "$(config_escape "$admin_token")"
    printf 'header = "Content-Type: application/json"\n'
    printf 'data-binary = "%s"\n' "$(config_escape "$payload")"
  } | curl --config - --silent --show-error --output /dev/null --write-out '%{http_code}' \
      --connect-timeout 10 --max-time 30 --max-redirs 0
)"; then
  die 'Identity provider request failed. Check connectivity and server logs.'
fi
unset admin_token client_secret payload
[[ "$status" == '201' ]] || die "Identity provider creation failed (HTTP $status). Check admin authorization, duplicate issuer/reference and server logs."
printf "Created identity provider '%s' (HTTP 201).\n" "$reference_name"
printf 'Register this Web redirect URI in the provider: %s\n' "$redirect_uri"
