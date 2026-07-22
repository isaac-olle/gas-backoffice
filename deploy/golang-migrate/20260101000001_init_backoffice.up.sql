CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE SCHEMA IF NOT EXISTS backoffice;

CREATE TABLE backoffice.users (
    id                     TEXT        NOT NULL PRIMARY KEY,
    admin_id               TEXT        NOT NULL,
    name                   TEXT        NOT NULL,
    role                   TEXT        NOT NULL DEFAULT 'tenant',
    is_system              BOOLEAN     NOT NULL DEFAULT FALSE,
    client_ids             TEXT[]      NOT NULL DEFAULT '{}',
    protected_resource_ids TEXT[]      NOT NULL DEFAULT '{}',
    author_id              TEXT        NOT NULL DEFAULT '',
    author_role            TEXT        NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE backoffice.protected_resources (
    id         TEXT        NOT NULL PRIMARY KEY,
    name       TEXT        NOT NULL,
    is_system  BOOLEAN     NOT NULL DEFAULT FALSE,
    user_id    TEXT        NOT NULL,
    scope_ids  TEXT[]      NOT NULL DEFAULT '{}',
    author_id  TEXT        NOT NULL DEFAULT '',
    author_role TEXT       NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.protected_resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE backoffice.scopes (
    id                    TEXT        NOT NULL PRIMARY KEY,
    name                  TEXT        NOT NULL,
    is_system             BOOLEAN     NOT NULL DEFAULT FALSE,
    protected_resource_id TEXT        NOT NULL,
    user_id               TEXT        NOT NULL,
    allowed_client_ids    TEXT[]      NOT NULL DEFAULT '{}',
    author_id             TEXT        NOT NULL DEFAULT '',
    author_role           TEXT        NOT NULL DEFAULT '',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.scopes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE UNIQUE INDEX scopes_name_protected_resource_id_unique
    ON backoffice.scopes (protected_resource_id, name)
    WHERE deleted_at IS NULL;

CREATE TABLE backoffice.clients (
    id             TEXT        NOT NULL PRIMARY KEY,
    name           TEXT        NOT NULL,
    user_id        TEXT        NOT NULL,
    allowed_flows  TEXT[]      NOT NULL DEFAULT '{}',
    allowed_scopes TEXT[]      NOT NULL DEFAULT '{}',
    redirect_uris  TEXT[]      NOT NULL DEFAULT '{}',
    key_ids        TEXT[]      NOT NULL DEFAULT '{}',
    secret_ids     TEXT[]      NOT NULL DEFAULT '{}',
    author_id      TEXT        NOT NULL DEFAULT '',
    author_role    TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.clients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

GRANT ALL ON SCHEMA backoffice TO gas_backoffice_admin;
GRANT ALL ON ALL TABLES IN SCHEMA backoffice TO gas_backoffice_admin;
GRANT USAGE ON SCHEMA backoffice TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA backoffice TO gas_backoffice_service;
