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

JAVA_HOME=$(/usr/libexec/java_home -v 21) ijhttp seed.http -e local -v http-client.env.json -p http-client.private.env.json
