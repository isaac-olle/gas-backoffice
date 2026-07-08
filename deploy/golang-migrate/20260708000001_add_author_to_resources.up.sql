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
