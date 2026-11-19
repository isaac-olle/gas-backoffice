# Register an OpenID identity provider

Use the administration API after the system has been seeded. This script calls
`POST /v1/backoffice/identity-providers`, so the application's creation service
encrypts the client secret with its configured AES key before saving it. Do not
pre-encrypt the secret or insert it directly into PostgreSQL.

Requires Bash, curl, jq and a GAS access token belonging to a
`system_administrator`. An ID token, Microsoft token, client-credentials token,
or `SEED_TOKEN` is not a substitute for that administrator access token.

## Microsoft Entra ID

In the Microsoft app registration, obtain the Directory (tenant) ID, Application
(client) ID, and a client secret **value** (not the secret ID). Configure a Web
redirect URI matching the GAS callback exactly.

From the repository root:

Fill in `scripts/identity-providers/microsoft.env` (gitignored), then run:

```bash
bash scripts/identity-providers/create.sh \
    --env-file scripts/identity-providers/microsoft.env \
    --api-base-url 'https://your-gas-api.example'
```

The Microsoft env file has exactly the same fields as `scripts/google_idp.env`,
replacing the GOOGLE_ prefix with MICROSOFT_ and retaining BACKOFFICE_REDIRECT_URI.
Fill in MICROSOFT_CLIENT_ID, MICROSOFT_CLIENT_SECRET, MICROSOFT_ISSUER,
MICROSOFT_AUTHORIZATION_ENDPOINT, MICROSOFT_TOKEN_ENDPOINT,
MICROSOFT_USERINFO_ENDPOINT, MICROSOFT_JWKS_URI and BACKOFFICE_REDIRECT_URI.
These endpoints come from your tenant-specific OpenID metadata document.

MICROSOFT_IDP_ID is included for parity with the seed configuration. The current
administration POST generates its own ID, so this field is not sent to the API.
The script prompts for the administrator access token, or reads
GAS_ADMIN_ACCESS_TOKEN from the environment. The env file uses Bash syntax and is
sourced; only load a trusted local file. Command-line options override its values.

Alternatively, provide non-secret settings on the command line:

```bash
bash scripts/identity-providers/create.sh \
    --api-base-url 'https://your-gas-api.example' \
    --tenant-id '<directory-tenant-guid>' \
    --client-id '<microsoft-application-client-id>'
```

The script prompts for the administrator access token and Microsoft client secret
without echoing them. It reads issuer, authorization, token, UserInfo and JWKS
endpoints from Microsoft's tenant-specific OpenID Discovery document. Defaults:

- Reference: `microsoft`; display name: `Microsoft`.
- Scopes: `openid`, `email`, `profile`.
- Redirect URI: `https://your-gas-api.example/v1/oidc/callback/microsoft`.

Override `--redirect-uri`, `--reference-name`, `--name` or `--scopes` if needed. Use a
concrete tenant GUID; `common`/`organizations` discovery has a templated issuer
that does not match GAS's exact issuer validation. This registers a tenant-
specific provider; it does not enable generic multi-tenant Microsoft login.

[Microsoft OpenID Connect documentation](https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols-oidc).

## Other OpenID providers

```bash
bash scripts/identity-providers/create.sh \
    --api-base-url 'https://your-gas-api.example' \
    --discovery-url 'https://issuer.example/.well-known/openid-configuration' \
    --reference-name 'example' --name 'Example' --client-id '<provider-client-id>'
```

Secrets can alternatively be supplied as environment variables
`GAS_ADMIN_ACCESS_TOKEN` and `OIDC_CLIENT_SECRET`. Non-secret defaults can be
provided as `API_BASE_URL`, `OIDC_TENANT_ID` and `OIDC_CLIENT_ID`.
On Windows run the script in Git Bash or WSL, with curl and jq installed.
The script does not print or save the request body, secret, or access token and
does not need `ENCRYPTION_KEY`. Use HTTPS for a remote GAS API.

HTTP 201 indicates creation succeeded. The endpoint generates the provider's
internal ID; select it at login using `idp_reference_name=microsoft` (or your
chosen reference). Errors are reported without dumping response bodies. Retrying
creation does not update an existing provider and can fail on unique constraints.
