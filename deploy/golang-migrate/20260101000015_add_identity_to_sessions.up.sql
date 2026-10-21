ALTER TABLE protocol.authorization_code_sessions
    ADD COLUMN IF NOT EXISTS idp_id              TEXT NOT NULL DEFAULT '' REFERENCES backoffice.identity_providers(id),
    ADD COLUMN IF NOT EXISTS id_token            TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS access_token        TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS refresh_token       TEXT,
    ADD COLUMN IF NOT EXISTS claims              JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS identity_expires_at TIMESTAMPTZ;
