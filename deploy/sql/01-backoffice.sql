
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE admins (
    id         TEXT        NOT NULL PRIMARY KEY,
    name       TEXT        NOT NULL,
    is_system  BOOLEAN     NOT NULL DEFAULT FALSE,
    user_ids TEXT[]      NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON admins
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE users (
    id                     TEXT        NOT NULL PRIMARY KEY,
    name                   TEXT        NOT NULL,
    is_system              BOOLEAN     NOT NULL DEFAULT FALSE,
    admin_id               TEXT        NOT NULL,
    client_ids             TEXT[]      NOT NULL DEFAULT '{}',
    protected_resource_ids TEXT[]      NOT NULL DEFAULT '{}',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE protected_resources (
    id         TEXT        NOT NULL PRIMARY KEY,
    name       TEXT        NOT NULL,
    is_system  BOOLEAN     NOT NULL DEFAULT FALSE,
    user_id  TEXT        NOT NULL,
    admin_id   TEXT        NOT NULL,
    scope_ids  TEXT[]      NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON protected_resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE scopes (
    id                    TEXT        NOT NULL PRIMARY KEY,
    name                  TEXT        NOT NULL,
    is_system             BOOLEAN     NOT NULL DEFAULT FALSE,
    protected_resource_id TEXT        NOT NULL,
    user_id             TEXT        NOT NULL,
    admin_id              TEXT        NOT NULL,
    allowed_client_ids    TEXT[]      NOT NULL DEFAULT '{}',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON scopes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE clients (
    id             TEXT        NOT NULL PRIMARY KEY,
    name           TEXT        NOT NULL,
    admin_id       TEXT        NOT NULL,
    user_id      TEXT        NOT NULL,
    allowed_flows  TEXT[]      NOT NULL DEFAULT '{}',
    allowed_scopes TEXT[]      NOT NULL DEFAULT '{}',
    redirect_uris  TEXT[]      NOT NULL DEFAULT '{}',
    key_ids        TEXT[]      NOT NULL DEFAULT '{}',
    secret_ids     TEXT[]      NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON clients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE admins             TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE users            TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE protected_resources TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE scopes             TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE clients            TO gas_backoffice;
