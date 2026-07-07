CREATE TABLE key_pairs (
    id                    TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    protected_resource_id TEXT        NOT NULL,
    -- data
    algorithm             TEXT        NOT NULL,
    public_key            BYTEA       NOT NULL,
    private_key           BYTEA       NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON key_pairs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE admin_assertion_keys (
    id          TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    admin_id    TEXT        NOT NULL,
    -- data
    algorithm   TEXT        NOT NULL,
    public_key  BYTEA       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON admin_assertion_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE client_assertion_keys (
    id          TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    client_id   TEXT        NOT NULL,
    -- data
    algorithm   TEXT        NOT NULL,
    public_key  BYTEA       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON client_assertion_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TABLE client_secrets (
    id          TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    client_id   TEXT        NOT NULL,
    -- data
    secret_hash TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON client_secrets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE key_pairs             TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE admin_assertion_keys  TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE client_assertion_keys TO gas_backoffice;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE client_secrets        TO gas_backoffice;
