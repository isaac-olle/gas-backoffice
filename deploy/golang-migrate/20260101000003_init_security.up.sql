CREATE SCHEMA security;

CREATE TABLE security.key_pairs (
    id                    TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    protected_resource_id TEXT        NOT NULL,
    -- data
    algorithm             TEXT        NOT NULL,
    public_key            BYTEA       NOT NULL,
    private_key           BYTEA       NOT NULL,
    author_id             TEXT        NOT NULL DEFAULT '',
    author_role           TEXT        NOT NULL DEFAULT '',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON security.key_pairs FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();

CREATE TABLE security.client_assertion_keys (
    id              TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    client_id       TEXT        NOT NULL,
    subscription_id TEXT        NOT NULL REFERENCES backoffice.subscriptions(id),
    -- data
    algorithm   TEXT        NOT NULL,
    public_key  BYTEA       NOT NULL,
    author_id   TEXT        NOT NULL DEFAULT '',
    author_role TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON security.client_assertion_keys FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();
CREATE INDEX client_assertion_keys_subscription_id_idx ON security.client_assertion_keys (subscription_id);

CREATE TABLE security.client_secrets (
    id              TEXT        NOT NULL PRIMARY KEY,
    -- ownership
    client_id       TEXT        NOT NULL,
    subscription_id TEXT        NOT NULL REFERENCES backoffice.subscriptions(id),
    -- data
    secret_hash TEXT        NOT NULL,
    author_id   TEXT        NOT NULL DEFAULT '',
    author_role TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON security.client_secrets FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();
CREATE INDEX client_secrets_subscription_id_idx ON security.client_secrets (subscription_id);
