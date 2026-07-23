CREATE SCHEMA IF NOT EXISTS events;

CREATE TABLE events.audit_events (
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

CREATE INDEX idx__audit_events__aggregate_id ON events.audit_events(aggregate_id);
CREATE INDEX idx__audit_events__event_type   ON events.audit_events(event_type);
CREATE INDEX idx__audit_events__occurred_at  ON events.audit_events(occurred_at);
