ALTER TABLE protocol.authorization_code_sessions
    ADD COLUMN occ_version BIGINT NOT NULL DEFAULT 1;
