ALTER TABLE clients
    ADD COLUMN key_ids    TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN secret_ids TEXT[] NOT NULL DEFAULT '{}';
