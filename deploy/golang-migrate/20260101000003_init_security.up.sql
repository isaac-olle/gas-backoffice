CREATE SCHEMA IF NOT EXISTS security;

CREATE TABLE IF NOT EXISTS security.key_pairs (
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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'security.key_pairs'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON security.key_pairs FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS security.client_assertion_keys (
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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'security.client_assertion_keys'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON security.client_assertion_keys FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;
CREATE INDEX IF NOT EXISTS client_assertion_keys_subscription_id_idx ON security.client_assertion_keys (subscription_id);

CREATE TABLE IF NOT EXISTS security.client_secrets (
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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'security.client_secrets'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON security.client_secrets FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;
CREATE INDEX IF NOT EXISTS client_secrets_subscription_id_idx ON security.client_secrets (subscription_id);


GRANT USAGE ON SCHEMA security TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA security TO gas_backoffice_service;
