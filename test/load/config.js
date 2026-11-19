export function required(name) {
  const value = __ENV[name];
  if (!value) throw new Error(`Missing ${name}; see test/load/README.md`);
  return value;
}

export function positiveInteger(name, fallback) {
  const value = Number(__ENV[name] || fallback);
  if (!Number.isInteger(value) || value < 1) {
    throw new Error(`${name} must be a positive integer`);
  }
  return value;
}

export const baseURL = required('BASE_URL').replace(/\/+$/, '');
if (!/^https?:\/\/[^/]+/.test(baseURL)) throw new Error('BASE_URL must be an absolute HTTP(S) URL');

export const selected = (__ENV.SCENARIOS || 'authorize,client_credentials').split(',').map((s) => s.trim());
const known = ['authorize', 'client_credentials', 'authorization_code'];
if (new Set(selected).size !== selected.length || selected.some((s) => !known.includes(s))) {
  throw new Error(`SCENARIOS must select distinct values from ${known.join(',')}`);
}

const profile = __ENV.PROFILE || 'smoke';
if (!['smoke', 'load', 'stress'].includes(profile)) throw new Error('PROFILE must be smoke, load or stress');

export function arrivalScenario(name) {
  const base = positiveInteger('BASE_RATE', 5);
  const peak = positiveInteger('PEAK_RATE', 50);
  if (peak < base) throw new Error('PEAK_RATE must be >= BASE_RATE');
  const common = {
    exec: name,
    timeUnit: '1s',
    preAllocatedVUs: positiveInteger('PREALLOCATED_VUS', 20),
    maxVUs: positiveInteger('MAX_VUS', 100),
    gracefulStop: '15s',
    tags: { endpoint: name === 'authorize' ? 'authorize' : 'token', grant: name },
  };
  if (common.maxVUs < common.preAllocatedVUs) throw new Error('MAX_VUS must be >= PREALLOCATED_VUS');
  if (profile !== 'stress') {
    return { ...common, executor: 'constant-arrival-rate',
      rate: positiveInteger('RATE', profile === 'smoke' ? 1 : 10),
      duration: __ENV.DURATION || (profile === 'smoke' ? '10s' : '1m') };
  }
  const ramp = __ENV.RAMP_DURATION || '30s';
  const hold = __ENV.HOLD_DURATION || '1m';
  return { ...common, executor: 'ramping-arrival-rate', startRate: 1,
    stages: [
      { duration: ramp, target: base }, { duration: hold, target: base },
      { duration: ramp, target: Math.ceil((base + peak) / 2) },
      { duration: hold, target: Math.ceil((base + peak) / 2) },
      { duration: ramp, target: peak }, { duration: hold, target: peak },
      { duration: '15s', target: 0 },
    ] };
}

export function thresholdsFor(names) {
  const thresholds = {};
  for (const name of names) {
    thresholds[`http_req_duration{scenario:${name}}`] = [
      `p(95)<${positiveInteger('P95_MS', 500)}`,
      `p(99)<${positiveInteger('P99_MS', 1500)}`,
    ];
    thresholds[`http_req_failed{scenario:${name}}`] = ['rate<0.01'];
    thresholds[`checks{scenario:${name}}`] = ['rate>=0.99'];
    if (name !== 'authorization_code') thresholds[`dropped_iterations{scenario:${name}}`] = ['count==0'];
  }
  return thresholds;
}
