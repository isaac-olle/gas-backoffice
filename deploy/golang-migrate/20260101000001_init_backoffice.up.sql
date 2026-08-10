CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE SCHEMA IF NOT EXISTS backoffice;

CREATE TABLE IF NOT EXISTS backoffice.users (
    id                     TEXT        NOT NULL PRIMARY KEY,
    admin_id               TEXT        NULLABLE,
    name                   TEXT        NOT NULL,
    email                  TEXT,
    idp_id                 TEXT,
    idp_subject            TEXT,
    role                   TEXT        NOT NULL DEFAULT 'tenant',
    is_system              BOOLEAN     NOT NULL DEFAULT FALSE,
    client_ids             TEXT[]      NOT NULL DEFAULT '{}',
    protected_resource_ids TEXT[]      NOT NULL DEFAULT '{}',
    author_id              TEXT        NOT NULL DEFAULT '',
    author_role            TEXT        NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ,

    CONSTRAINT users_idp_fk
    FOREIGN KEY (idp_id) REFERENCES backoffice.identity_providers(id),

    CONSTRAINT users_idp_identity_unique
    UNIQUE (idp_id, idp_subject)
    );
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'backoffice.users'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.users FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS backoffice.protected_resources (
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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'backoffice.protected_resources'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.protected_resources FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS backoffice.scopes (
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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'backoffice.scopes'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.scopes FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;
CREATE UNIQUE INDEX IF NOT EXISTS scopes_name_protected_resource_id_unique
    ON backoffice.scopes (protected_resource_id, name)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS backoffice.clients (
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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_updated_at' AND tgrelid = 'backoffice.clients'::regclass) THEN
    CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.clients FOR EACH ROW EXECUTE FUNCTION update_updated_at();
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS backoffice.sessions (
    id             TEXT        NOT NULL PRIMARY KEY,
    user_id        TEXT        NOT NULL,
    author_id      TEXT        NOT NULL DEFAULT '',
    author_role    TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    logged_out_at  TIMESTAMPTZ
);

GRANT USAGE ON SCHEMA backoffice TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA backoffice TO gas_backoffice_service;

CREATE TABLE IF NOT EXISTS backoffice.oauth_identity_providers (
    id                      TEXT        NOT NULL PRIMARY KEY,
    name                    TEXT        NOT NULL,
    issuer                  TEXT        NOT NULL,
    client_id               TEXT        NOT NULL,
    client_secret           TEXT,
    authorization_endpoint  TEXT        NOT NULL,
    token_endpoint          TEXT        NOT NULL,
    userinfo_endpoint       TEXT,
    jwks_uri                TEXT        NOT NULL,
    scopes                  TEXT[]      NOT NULL DEFAULT '{openid,email,profile}',
    enabled                 BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT identity_providers_issuer_unique UNIQUE (issuer)
);

DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_trigger
    WHERE tgname = 'set_updated_at'
      AND tgrelid = 'backoffice.identity_providers'::regclass
  ) THEN
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON backoffice.identity_providers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
END IF;
END $$;