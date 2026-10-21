-- Rename table to match repository expectations
ALTER TABLE protocol.protocol_sessions RENAME TO protocol_authorization_code_sessions;

-- Add columns that the repository reads and writes
ALTER TABLE protocol.protocol_authorization_code_sessions
    ADD COLUMN IF NOT EXISTS client_owner_user_id  TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS protected_resource_id TEXT NOT NULL DEFAULT '';
