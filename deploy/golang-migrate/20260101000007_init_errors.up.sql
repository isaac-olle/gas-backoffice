CREATE TABLE public.errors (
    error_id    TEXT NOT NULL PRIMARY KEY,
    type        TEXT NOT NULL,
    code        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    details     JSONB,
    cause       TEXT[] NOT NULL DEFAULT '{}',
    stack_trace JSONB,
    request_id  TEXT,
    occurred_at TIMESTAMPTZ
);
