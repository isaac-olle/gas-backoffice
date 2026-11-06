CREATE SCHEMA backoffice;

CREATE FUNCTION backoffice.update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE backoffice.identity_providers (
    id                      TEXT        NOT NULL PRIMARY KEY,
    reference_name          TEXT        NOT NULL UNIQUE,
    name                    TEXT        NOT NULL,
    issuer                  TEXT        NOT NULL,
    client_id               TEXT        NOT NULL,
    client_secret           TEXT,
    authorization_endpoint  TEXT        NOT NULL,
    token_endpoint          TEXT        NOT NULL,
    userinfo_endpoint       TEXT,
    jwks_uri                TEXT        NOT NULL,
    scopes                  TEXT[]      NOT NULL DEFAULT '{openid,email,profile}',
    backoffice_redirect_uri TEXT,
    enabled                 BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT identity_providers_issuer_unique UNIQUE (issuer)
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.identity_providers FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();

CREATE TABLE backoffice.users (
    id                     TEXT        NOT NULL PRIMARY KEY,
    name                   TEXT        NOT NULL,
    email                  TEXT,
    system_role            TEXT        NOT NULL DEFAULT 'subscription_user',
    is_system              BOOLEAN     NOT NULL DEFAULT FALSE,
    subscription_ids       TEXT[]      NOT NULL DEFAULT '{}',
    memberships            JSONB       NOT NULL DEFAULT '[]'::jsonb,
    version                INTEGER     NOT NULL DEFAULT 1,
    author_id              TEXT        NOT NULL DEFAULT '',
    author_role            TEXT        NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
    );
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.users FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();

CREATE TABLE backoffice.subscriptions (
    id         TEXT        NOT NULL PRIMARY KEY,
    name       TEXT        NOT NULL,
    owner_id   TEXT        NOT NULL REFERENCES backoffice.users(id),
    is_system  BOOLEAN     NOT NULL DEFAULT FALSE,
    version    INTEGER     NOT NULL DEFAULT 1,
    author_id  TEXT        NOT NULL DEFAULT '',
    author_role TEXT       NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.subscriptions FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();

CREATE TABLE backoffice.subscription_members (
    subscription_id TEXT NOT NULL REFERENCES backoffice.subscriptions(id) ON DELETE CASCADE,
    user_id         TEXT NOT NULL REFERENCES backoffice.users(id) ON DELETE CASCADE,
    subscription_role TEXT NOT NULL DEFAULT 'user' CHECK (subscription_role IN ('user')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (subscription_id, user_id)
);
CREATE INDEX subscription_members_user_id_idx ON backoffice.subscription_members (user_id);

CREATE TABLE backoffice.protected_resources (
    id              TEXT        NOT NULL PRIMARY KEY,
    name            TEXT        NOT NULL,
    is_system       BOOLEAN     NOT NULL DEFAULT FALSE,
    enabled         BOOLEAN     NOT NULL DEFAULT TRUE,
    subscription_id TEXT        NOT NULL REFERENCES backoffice.subscriptions(id),
    scope_ids       TEXT[]      NOT NULL DEFAULT '{}',
    version         INTEGER     NOT NULL DEFAULT 1,
    author_id       TEXT        NOT NULL DEFAULT '',
    author_role     TEXT        NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.protected_resources FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();
CREATE INDEX protected_resources_subscription_id_idx ON backoffice.protected_resources (subscription_id);
CREATE UNIQUE INDEX protected_resources_subscription_id_name_unique
    ON backoffice.protected_resources (subscription_id, name)
    WHERE deleted_at IS NULL;

CREATE TABLE backoffice.scopes (
    id                    TEXT        NOT NULL PRIMARY KEY,
    name                  TEXT        NOT NULL,
    is_system             BOOLEAN     NOT NULL DEFAULT FALSE,
    enabled               BOOLEAN     NOT NULL DEFAULT TRUE,
    protected_resource_id TEXT        NOT NULL,
    subscription_id       TEXT        NOT NULL REFERENCES backoffice.subscriptions(id),
    allowed_client_ids    TEXT[]      NOT NULL DEFAULT '{}',
    version               INTEGER     NOT NULL DEFAULT 1,
    author_id             TEXT        NOT NULL DEFAULT '',
    author_role           TEXT        NOT NULL DEFAULT '',
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.scopes FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();
CREATE UNIQUE INDEX scopes_name_protected_resource_id_unique
    ON backoffice.scopes (protected_resource_id, name)
    WHERE deleted_at IS NULL;
CREATE INDEX scopes_subscription_id_idx ON backoffice.scopes (subscription_id);

CREATE TABLE backoffice.clients (
    id              TEXT        NOT NULL PRIMARY KEY,
    name            TEXT        NOT NULL,
    subscription_id TEXT        NOT NULL REFERENCES backoffice.subscriptions(id),
    is_system       BOOLEAN     NOT NULL DEFAULT FALSE,
    enabled        BOOLEAN     NOT NULL DEFAULT TRUE,
    public         BOOLEAN     NOT NULL DEFAULT FALSE,
    allowed_grant_types TEXT[] NOT NULL DEFAULT '{}',
    redirect_uris  TEXT[]      NOT NULL DEFAULT '{}',
    assertion_key_ids        TEXT[]      NOT NULL DEFAULT '{}',
    secret_ids     TEXT[]      NOT NULL DEFAULT '{}',
    allowed_scopes TEXT[]      NOT NULL DEFAULT '{}',
    version       INTEGER     NOT NULL DEFAULT 1,
    author_id      TEXT        NOT NULL DEFAULT '',
    author_role    TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON backoffice.clients FOR EACH ROW EXECUTE FUNCTION backoffice.update_updated_at();
CREATE INDEX clients_subscription_id_idx ON backoffice.clients (subscription_id);
CREATE UNIQUE INDEX clients_subscription_id_name_unique
    ON backoffice.clients (subscription_id, name)
    WHERE deleted_at IS NULL;

CREATE TABLE backoffice.sessions (
    id             TEXT        NOT NULL PRIMARY KEY,
    user_id        TEXT        NOT NULL,
    user_role      TEXT        NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    logged_out_at  TIMESTAMPTZ
);
