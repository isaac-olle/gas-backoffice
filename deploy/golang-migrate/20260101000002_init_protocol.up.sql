CREATE SCHEMA protocol;

CREATE TABLE protocol.authorization_code_sessions (
    session_id     TEXT        NOT NULL PRIMARY KEY,
    client_id      TEXT        NOT NULL,
    redirect_uri   TEXT        NOT NULL,
    scope          TEXT        NOT NULL,
    state          TEXT        NOT NULL DEFAULT '',
    response_types JSONB       NOT NULL DEFAULT '[]',
    code           TEXT        NOT NULL DEFAULT '',
    pkce           JSONB,
    session_state  TEXT        NOT NULL DEFAULT '',
    grant_type     TEXT        NOT NULL DEFAULT '',
    client_owner_user_id  TEXT NOT NULL DEFAULT '',
    protected_resource_id TEXT NOT NULL DEFAULT '',
    idp_id         TEXT        NOT NULL DEFAULT '' REFERENCES backoffice.identity_providers(id),
    id_token       TEXT        NOT NULL DEFAULT '',
    access_token   TEXT        NOT NULL DEFAULT '',
    refresh_token  TEXT,
    claims         JSONB       NOT NULL DEFAULT '{}',
    identity_expires_at TIMESTAMPTZ,
    consumed       BOOLEAN     NOT NULL DEFAULT FALSE,
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON protocol.authorization_code_sessions FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();

CREATE TABLE protocol.client_credentials_sessions (
    id         TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    client_id  TEXT        NOT NULL,
    user_id    TEXT        NOT NULL,
    -- data
    scope      TEXT        NOT NULL,
    issued_at  TIMESTAMPTZ NOT NULL,
    protected_resource_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
