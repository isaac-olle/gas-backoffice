ALTER TABLE users
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE protected_resources
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE clients
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE scopes
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE key_pairs
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE admin_assertion_keys
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE client_assertion_keys
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

ALTER TABLE client_secrets
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';
