ALTER TABLE protocol.authorization_code_sessions
    ADD COLUMN client_owner_user_id TEXT NOT NULL DEFAULT '';
