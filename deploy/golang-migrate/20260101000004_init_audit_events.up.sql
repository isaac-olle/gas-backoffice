CREATE SCHEMA IF NOT EXISTS events;

CREATE TABLE IF NOT EXISTS events.audit_events (
    event_id       UUID        NOT NULL PRIMARY KEY,
    event_type     TEXT        NOT NULL,
    aggregate_id   TEXT        NOT NULL,
    aggregate_type TEXT        NOT NULL,
    author_id      TEXT,
    author_role    TEXT,
    payload        JSONB       NOT NULL,
    occurred_at    TIMESTAMPTZ NOT NULL,
    recorded_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx__audit_events__aggregate_id ON events.audit_events(aggregate_id);
CREATE INDEX IF NOT EXISTS idx__audit_events__event_type   ON events.audit_events(event_type);
CREATE INDEX IF NOT EXISTS idx__audit_events__occurred_at  ON events.audit_events(occurred_at);

GRANT USAGE ON SCHEMA events TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA events TO gas_backoffice_service;
