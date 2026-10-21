CREATE TABLE IF NOT EXISTS public.errors (
    error_id    TEXT NOT NULL PRIMARY KEY,
    type        TEXT NOT NULL,
    code        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    details     JSONB,
    cause       TEXT[] NOT NULL DEFAULT '{}',
    stack_trace JSONB
);

GRANT SELECT, INSERT, UPDATE, DELETE ON public.errors TO gas_backoffice_service;
