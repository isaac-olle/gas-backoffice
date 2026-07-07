
CREATE TABLE public.audit_events (
    event_id       UUID        NOT NULL PRIMARY KEY,
    event_type     TEXT        NOT NULL,
    aggregate_id   TEXT        NOT NULL,
    aggregate_type TEXT        NOT NULL,
    actor_id       TEXT,
    actor_role     TEXT,
    payload        JSONB       NOT NULL,
    occurred_at    TIMESTAMPTZ NOT NULL,
    recorded_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx__audit_events__aggregate_id ON public.audit_events(aggregate_id);
CREATE INDEX idx__audit_events__event_type   ON public.audit_events(event_type);
CREATE INDEX idx__audit_events__occurred_at  ON public.audit_events(occurred_at);

GRANT SELECT, INSERT ON TABLE public.audit_events TO gas_backoffice;
