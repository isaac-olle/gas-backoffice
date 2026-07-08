
CREATE TABLE protocol_sessions (
    session_id     TEXT        NOT NULL PRIMARY KEY,
    client_id      TEXT        NOT NULL,
    redirect_uri   TEXT        NOT NULL,
    scope          TEXT        NOT NULL,
    state          TEXT        NOT NULL DEFAULT '',
    response_types JSONB       NOT NULL DEFAULT '[]',
    code           TEXT        NOT NULL DEFAULT '',
    pkce           JSONB,
    session_state  TEXT        NOT NULL DEFAULT '',
    flow           TEXT        NOT NULL DEFAULT '',
    consumed       BOOLEAN     NOT NULL DEFAULT FALSE,
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON protocol_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE protocol_client_credentials_sessions (
    session_id  TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    client_id   TEXT        NOT NULL,
    user_id   TEXT        NOT NULL,
    admin_id    TEXT        NOT NULL,
    -- data
    scope       TEXT        NOT NULL,
    issued_at   TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE protocol_sessions                    TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE protocol_client_credentials_sessions TO gas_backoffice;
