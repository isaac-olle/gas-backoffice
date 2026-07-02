\connect gas

CREATE TABLE public.outbox_event_bus (
    event_id     UUID        NOT NULL PRIMARY KEY,
    event_name   TEXT        NOT NULL,
    aggregate_id TEXT        NOT NULL,
    payload      JSONB       NOT NULL,
    actor_id     TEXT,
    actor_role   TEXT,
    occurred_at  TIMESTAMPTZ NOT NULL,
    inserted_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx__outbox_event_bus__inserted_at ON public.outbox_event_bus(inserted_at);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.outbox_event_bus TO gas_backoffice;
