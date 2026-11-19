# Protocol load tests (k6)

Targets the real GAS endpoints, with a preconfigured test database and running
API. This does not start the API, create clients, or seed credentials. Use a
dedicated load-test client/resource and a disposable database: successful token
requests create sessions and audit/outbox records.

## Scenarios

| SCENARIOS value | Request | Successful result |
| --- | --- | --- |
| authorize | GET /oidc/authorize with fresh state and S256 PKCE | 302 to the configured IdP authorization endpoint |
| client_credentials | POST /oidc/token with Basic authentication and explicit scopes | 200 with access token and cache-prevention headers |
| authorization_code | POST /oidc/token with a fresh code and matching verifier | 200 with access token and cache-prevention headers |

Select one or comma-separated scenarios. The default is
`authorize,client_credentials`. These run concurrently, with independent rates.
For example, a peak of 50 means up to 50 authorize and 50 token requests per
second when both arrival-rate scenarios are selected.

The authorize test stops at the 302. It does not follow redirects or send load
to Google/the external IdP, and it does not create authorization codes. A 302
to an error callback fails its checks. The token test exercises actual credential
verification, authorization, persistence and signing, rather than only invalid
requests.

## Run a smoke test first

Install k6, then set the environment variables listed in `.env.example`.
k6 does not automatically load that file. Example in PowerShell:

```powershell
$env:BASE_URL = 'http://localhost:8090'
$env:CLIENT_ID = '<test-client-id>'
$env:CLIENT_SECRET = '<test-client-secret>'
$env:REDIRECT_URI = 'http://localhost:3000/callback'
$env:IDP_REFERENCE_NAME = '<configured-test-idp-reference>'
$env:IDP_AUTHORIZATION_URL = 'https://test-idp.example/authorize'
$env:AUTHORIZE_SCOPE = 'openid'
$env:TOKEN_SCOPE = 'read'
k6 run test/load/protocol.js
```

The client must allow the selected grants, have the redirect URI registered,
and be authorized for the requested scopes. Scopes in each request must belong
to one protected resource. The configured IdP must be enabled. Use the actual
configured IdP authorization endpoint for `IDP_AUTHORIZATION_URL`.

The smoke profile sends 1 request/second per scenario for 10 seconds. A public
client can test authorize and authorization-code exchange without CLIENT_SECRET;
client_credentials requires a confidential client and its secret.

For only token issuance:

```powershell
$env:SCENARIOS = 'client_credentials'
k6 run test/load/protocol.js
```

## Load and stress profiles

```sh
# With configuration already exported in the environment:
k6 run -e PROFILE=load test/load/protocol.js
k6 run -e PROFILE=stress -e BASE_RATE=5 -e PEAK_RATE=50 test/load/protocol.js
# Equivalent Make targets:
make load-test
make stress-test
```

`make load-test` uses the configured PROFILE, defaulting to smoke.

| Profile | Default per-scenario traffic |
| --- | --- |
| smoke | 1 request/s for 10s |
| load | 10 requests/s for 1m |
| stress | Ramp 1 -> 5 -> 28 -> 50 requests/s, holding each level for 1m, then ramp down |

Override RATE/DURATION for smoke or load; BASE_RATE, PEAK_RATE,
RAMP_DURATION and HOLD_DURATION for stress. PREALLOCATED_VUS and MAX_VUS control
load-generator concurrency. Arrival-rate executors maintain the requested
traffic independently of server latency; dropped iterations fail the run so a
saturated load generator cannot quietly report a successful test.

Default thresholds, evaluated separately per scenario:

- HTTP failures below 1% and successful checks at least 99%.
- p95 latency below 500ms; p99 below 1500ms (override P95_MS/P99_MS).
- No dropped arrival-rate iterations.

These are initial targets, not an established service SLA. Compare request rate,
latency percentiles, errors, CPU, database connections and outbox backlog to find
the saturation point. Running on a shared CI runner is not a stable performance
baseline; these tests are opt-in and are not attached to every deployment.

## Authorization-code exchange

Authorization codes are single-use. Provide a JSON array of freshly issued GAS
codes in `test/load/codes.local.json`, following `codes.example.json`. Each code
must belong to CLIENT_ID and REDIRECT_URI and include its matching code_verifier
when issued with PKCE. Codes must remain valid throughout the run (currently
15 minutes after issuance). Never reuse a fixture after a run.

```powershell
$env:SCENARIOS = 'authorization_code'
$env:AUTH_CODES_FILE = './codes.local.json'
$env:CODE_VUS = '10'
k6 run test/load/protocol.js
```

The fixture path is relative to protocol.js. One iteration consumes one unique
code using k6's scenario-wide iteration index, even with concurrent VUs.
This scenario uses shared iterations, not the arrival-rate profile: its total
requests equal the fixture length. If CODE_MAX_DURATION (default 5m) expires
before all codes are consumed, the iteration-count threshold fails.

The suite does not automate external IdP login or provision this code batch.
For sustained authorization-code load, generate codes through a local test IdP
and the normal GAS callback; do not replay a single code or load-test Google.

`.local.json` fixtures and `results/` are gitignored. Scripts do not log tokens,
credentials or request bodies, and metric labels use fixed endpoint names.
Avoid k6 HTTP debug logging when using real test credentials.

Reference: [k6 arrival-rate executors](https://grafana.com/docs/k6/latest/using-k6/scenarios/executors/ramping-arrival-rate/).
