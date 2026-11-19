import http from 'k6/http';
import { check } from 'k6';
import crypto from 'k6/crypto';
import encoding from 'k6/encoding';
import execution from 'k6/execution';
import { SharedArray } from 'k6/data';
import { arrivalScenario, baseURL, positiveInteger, required, selected, thresholdsFor } from './config.js';

const authorizeEnabled = selected.includes('authorize');
const ccEnabled = selected.includes('client_credentials');
const codeEnabled = selected.includes('authorization_code');
const clientID = required('CLIENT_ID');
const clientSecret = ccEnabled ? required('CLIENT_SECRET') : (__ENV.CLIENT_SECRET || '');
const redirectURI = authorizeEnabled || codeEnabled ? required('REDIRECT_URI') : '';
const idp = authorizeEnabled ? required('IDP_REFERENCE_NAME') : '';
const idpEndpoint = authorizeEnabled ? required('IDP_AUTHORIZATION_URL') : '';
const authorizeScope = authorizeEnabled ? required('AUTHORIZE_SCOPE') : '';
const tokenScope = ccEnabled ? required('TOKEN_SCOPE') : '';

const codes = codeEnabled ? new SharedArray('authorization codes', () => {
  const rows = JSON.parse(open(required('AUTH_CODES_FILE')));
  if (!Array.isArray(rows) || rows.length === 0) throw new Error('AUTH_CODES_FILE must be a nonempty JSON array');
  const seen = new Set();
  for (const row of rows) {
    if (typeof row.code !== 'string' || !row.code || seen.has(row.code)) {
      throw new Error('Every authorization code must be a unique, nonempty string');
    }
    if (row.code_verifier !== undefined && typeof row.code_verifier !== 'string') {
      throw new Error('code_verifier must be a string');
    }
    seen.add(row.code);
  }
  return rows;
}) : [];

const scenarios = {};
for (const name of selected) {
  scenarios[name] = name === 'authorization_code' ? {
    executor: 'shared-iterations', exec: name,
    vus: Math.min(positiveInteger('CODE_VUS', 10), codes.length),
    iterations: codes.length, maxDuration: __ENV.CODE_MAX_DURATION || '5m',
    tags: { endpoint: 'token', grant: 'authorization_code' },
  } : arrivalScenario(name);
}

export const options = {
  scenarios,
  thresholds: thresholdsFor(selected),
  maxRedirects: 0,
  discardResponseBodies: true,
  // Never attach a fresh state/code to metric labels or emit tokens in logs.
  systemTags: ['status', 'method', 'name', 'scenario', 'expected_response', 'check'],
};

if (codeEnabled) {
  options.thresholds['iterations{scenario:authorization_code}'] = [`count==${codes.length}`];
}

function form(values) {
  return Object.keys(values).map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(values[key])}`).join('&');
}

function tokenRequest(values) {
  const headers = { 'Content-Type': 'application/x-www-form-urlencoded' };
  if (clientSecret) {
    headers.Authorization = `Basic ${encoding.b64encode(`${clientID}:${clientSecret}`)}`;
  } else {
    values.client_id = clientID; // Public authorization-code clients authenticate with PKCE.
  }
  const response = http.post(`${baseURL}/oidc/token`, form(values), {
    headers, redirects: 0, responseType: 'text', timeout: __ENV.REQUEST_TIMEOUT || '10s',
    tags: { name: 'POST /oidc/token' }, responseCallback: http.expectedStatuses(200),
  });
  let body;
  try { body = response.json(); } catch (_) { body = null; }
  check(response, {
    'token: HTTP 200': (r) => r.status === 200,
    'token: access token returned': () => body !== null && typeof body.access_token === 'string' && body.access_token.length > 0,
    'token: Bearer with positive expiry': () => body !== null && body.token_type === 'Bearer' && body.expires_in > 0,
    'token: no HTTP cache': (r) => r.headers['Cache-Control'] === 'no-store' && r.headers.Pragma === 'no-cache',
  });
}

export function authorize() {
  const verifier = encoding.b64encode(crypto.randomBytes(32), 'rawurl');
  const challenge = encoding.b64encode(crypto.sha256(verifier, 'binary'), 'rawurl');
  const query = form({
    response_type: 'code', client_id: clientID, redirect_uri: redirectURI,
    scope: authorizeScope, idp_reference_name: idp,
    state: encoding.b64encode(crypto.randomBytes(24), 'rawurl'),
    code_challenge: challenge, code_challenge_method: 'S256',
  });
  const response = http.get(`${baseURL}/oidc/authorize?${query}`, {
    redirects: 0, timeout: __ENV.REQUEST_TIMEOUT || '10s', tags: { name: 'GET /oidc/authorize' },
    responseCallback: http.expectedStatuses(302),
  });
  check(response, {
    'authorize: HTTP 302': (r) => r.status === 302,
    'authorize: redirect to configured IdP': (r) => {
      const location = r.headers.Location || '';
      return location.split('?')[0] === idpEndpoint.split('?')[0] && /[?&]state=/.test(location);
    },
  });
}

export function client_credentials() {
  tokenRequest({ grant_type: 'client_credentials', scope: tokenScope });
}

export function authorization_code() {
  // The scenario-wide index allocates each single-use code exactly once.
  const row = codes[execution.scenario.iterationInTest];
  tokenRequest({ grant_type: 'authorization_code', code: row.code,
    redirect_uri: redirectURI, code_verifier: row.code_verifier || '' });
}
