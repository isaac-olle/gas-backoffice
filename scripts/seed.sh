#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

set -a
source ../.env.local
source google_idp.env
set +a

JAVA_HOME=$(/usr/libexec/java_home -v 21) ijhttp seed.http -e local -v http-client.env.json -p http-client.private.env.json
